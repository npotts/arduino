/*
The MIT License (MIT)

Copyright (c) 2016 Nick Potts

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/gonum/plot/vg"
	"github.com/jmoiron/sqlx"
	"github.com/npotts/arduino/WxStation/plots"
	"github.com/npotts/homehub"
	"github.com/pkg/errors"
	"github.com/tarm/serial"
	"image/color"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" //mysql support
	_ "github.com/lib/pq"              //postgres support
	_ "github.com/mattn/go-sqlite3"    //sqlite3  support
)

var (
	app    = kingpin.New("WxStation", "Shovel data coming from an arduino configured as a WxStation to a brianiac instance")
	table  = app.Flag("table", "The database table to fire into (daemon mode) or read from (plot)").Default("wxstation").Short('t').String()
	daemon = app.Command("demon", "Operate in Daemon mode - shovel RS232 data to brainiac instance")
	baud   = daemon.Flag("baud", "The baud rate to listen at.  Default is the compiled in baud rate").Short('b').Default("115200").Int()
	device = daemon.Arg("device", "The RS232 serial device connected to the Arduino running WxStation (http://github.com/npotts/arduino/WxStation)").Required().String()
	url    = daemon.Arg("url", "URL to brainaic instance").Required().URL()

	plotc  = app.Command("plot", "Generate plots of data")
	driver = plotc.Flag("driver", "Database driver.  Supported drivers are mysql, postgresql, and sqlite3").Default("mysql").Short('d').String()
	dsn    = plotc.Flag("dsn", "DSN / Dial string").Default("brainiac:brainiac@/brainiac?parseTime=true&loc=UTC").Short('D').String()
	output = plotc.Arg("svg", "Path to SVG output file").Default("trh.svg").String()
)

func getPort() *serial.Port {
	fmt.Printf("Attempting to open %s: ", *device)
	port, err := serial.OpenPort(&serial.Config{Name: *device, Baud: *baud, Size: 8, Parity: serial.ParityNone, StopBits: 1})
	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to open serial device"))
		os.Exit(1)
	}
	fmt.Println("done")
	return port
}

func do(method string, line []byte) bool {
	x := homehub.Datam{}
	data := fmt.Sprintf(`{"table":%q,"data":%s}`, *table, strings.Replace(strings.Replace(string(line), "\n", "", -1), "\r", "", -1))
	err := json.Unmarshal([]byte(data), &x)
	if err != nil {
		return false
	}

	if req, err := http.NewRequest(method, (*url).String(), strings.NewReader(data)); err == nil {
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: ", resp)
			return false
		}
		return (resp.StatusCode == http.StatusOK)
	}
	return false
}

var first = true
var count = 0

func monitor() {
	port := getPort()
	rdr := bufio.NewReader(port)
	for {
		line, err := rdr.ReadBytes('\r')
		if err == nil {
			if first {
				first = !do("PUT", line)
			}
			if do("POST", line) {
				count++
				fmt.Printf("\r%d", count)
			}
		}
	}
}

func plot() {
	now := time.Now().UTC().Add(-2 * 24 * time.Hour)
	file, err := os.Create(*output)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	sq, err := sqlx.Open(*driver, *dsn)
	if err != nil {
		panic(err)
	}
	defer sq.Close()
	data := plots.Measurements{}
	if e := sq.Select(&data, "SELECT * FROM "+*table); e != nil {
		panic(e)
	}

	tp := plots.NewTimePlot("Title", "Previous Hours", "Temp / RH")

	for key, color := range map[string]color.RGBA{
		"battery": color.RGBA{R: 255, G: 255, B: 255, A: 255},
		// "humidity":      color.RGBA{R: 255, G: 255, B: 255, A: 255},
		// "humiditytemp":  color.RGBA{R: 255, G: 255, B: 255, A: 255},
		// "ihumidity":     color.RGBA{R: 255, G: 255, B: 255, A: 255},
		// "ihumiditytemp": color.RGBA{R: 255, G: 255, B: 255, A: 255},
		// "pressure":       color.RGBA{R: 255, G: 255, B: 255, A: 255},
		// "pressuretemp":   color.RGBA{R: 255, G: 255, B: 255, A: 255},
		// "temperature":    color.RGBA{R: 255, G: 255, B: 255, A: 255},
		// "temperatureext": color.RGBA{R: 255, G: 255, B: 255, A: 255},
		// "vref":           color.RGBA{R: 255, G: 255, B: 255, A: 255},
	} {
		tp.AddTrace(plots.Trace{Data: data.XYs(key, now), Color: color})
	}

	if err := tp.WriteTo(file, 5*vg.Inch, 3*vg.Inch, "svg"); err != nil {
		panic(err)
	}
}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case daemon.FullCommand():
		monitor()
	case plotc.FullCommand():
		plot()

	}

}
