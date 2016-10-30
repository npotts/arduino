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
	"github.com/npotts/homehub"
	"github.com/pkg/errors"
	"github.com/tarm/serial"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var (
	app    = kingpin.New("wxstation-proxy", "Shovel data from an arduino via serial port to a brianiac instance")
	baud   = app.Flag("baud", "The baud rate to listen at").Short('b').Default("115200").Int()
	device = app.Arg("device", "The RS232 serial device to read from").Required().String()
	table  = app.Arg("table", "The database table to fire into").Required().String()
	url    = app.Arg("url", "URL to brainaic instance").Required().URL()
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

func main() {
	app.Parse(os.Args[1:])
	monitor()
}
