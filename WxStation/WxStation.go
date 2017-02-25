/*
The MIT License (MIT)

Copyright (c) 2016-2017 Nick Potts

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
	"github.com/pkg/errors"
	"github.com/tarm/serial"
	"os"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

var (
	app      = kingpin.New("WxStation2", "Shovel data coming from an arduino configured as a WxStation into a InfluxDB instance")
	baud     = app.Flag("baud", "The baud rate to listen at.  Default is the compiled in baud rate").Short('b').Default("115200").Int()
	device   = app.Flag("device", "The RS232 serial device connected to the Arduino running WxStation (http://github.com/npotts/arduino/WxStation)").Default("/dev/ttyUSB0").String()
	db       = app.Flag("table", "The database to fire into").Short('t').Default("wx").String()
	user     = app.Flag("user", "The database table to fire into").Short('u').Default("wxstation").String()
	password = app.Flag("password", "The database table to fire into").Short('p').Default("wxstation").String()
	influxdb = app.Flag("influx", `URL to Influx DB instance.  Usually this is something like 'http://server:8086', but may be somewhere else.`).Default("http://pika:8086").String()
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

func getClient() client.Client {
	fmt.Printf("Attempting to connect to influx db instance: ")
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               *influxdb,
		Username:           *user,
		Password:           *password,
		UserAgent:          "WxStations2/arduino",
		Timeout:            0,
		InsecureSkipVerify: false,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return c
}

var count = 0

func writeLine(line []byte, c client.Client) error {
	m := map[string]interface{}{"count": count}

	if err := json.Unmarshal(line, &m); err != nil {
		return fmt.Errorf("Unpack error: %v: %v", string(line), err)
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "us",
		Database:  *db,
	})

	if err != nil {
		return fmt.Errorf("Unable to form BatchPoints: %v", err)
	}

	pt, err := client.NewPoint("sample", map[string]string{}, m) //use servers timestamp
	// pt, err := client.NewPoint("sample", map[string]string{}, m, time.Now())
	if err != nil {
		return fmt.Errorf("Unable to create point: %v", err)
	}

	bp.AddPoint(pt)
	count++
	return c.Write(bp)
}

func monitor() {

	port := getPort()
	defer port.Close()
	ic := getClient()
	defer ic.Close()
	rdr := bufio.NewReader(port)
	for {
		line, err := rdr.ReadBytes('\r')
		if err == nil {
			if e := writeLine(line, ic); e != nil {
				fmt.Println(e)
			}
		}
	}
}

func main() {
	app.Parse(os.Args[1:])
	// c := getClient()
	// fmt.Println(writeLine([]byte(`{"vref":0.00,"battery":4.37,"humidity":29,"humidityTemp":1,"pressure":842.40,"pressureTemp":-2.13,"ihumidity":55.05,"ihumidityTemp":0.09,"temperatureExt":-5.13,"temperature":-5.00}`), c))

	monitor()
}
