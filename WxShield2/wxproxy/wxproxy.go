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
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/tarm/serial"
	"os"
	"strings"
)

var (
	create = "CREATE TABLE IF NOT EXISTS %s (id SERIAL, timestamp TIMESTAMP DEFAULT now(), pressure FLOAT, tempa FLOAT, tempb FLOAT, humidity FLOAT, ptemp FLOAT, htemp FLOAT, battery FLOAT, indx INT)"

	app            = kingpin.New("wxProxy", "Shovel data from an arduino into a postgres/cockroachdb")
	dataSourceName = app.Flag("dataSource", "Where should we connect to and yelp at (usually a string like 'postgres://user:password@server/database?sslmode=disable')").Short('s').Default("postgres://wx:wx@pika/wx?sslmode=disable").String()
	table          = app.Flag("table", "The database table to fire into").Short('t').Default("raw").String()
	device         = app.Flag("device", "The RS232 serial device to read from (at 57600 baud)").Short('D').Default("/dev/ttyUSB0").String()
)

type inserter struct {
	db  *sqlx.DB
	ser *serial.Port
}

func getPort() *serial.Port {
	port, err := serial.OpenPort(&serial.Config{Name: *device, Baud: 57600, Size: 8, Parity: serial.ParityNone, StopBits: 1})
	if err != nil {
		panic(errors.Wrap(err, "Unable to open serial device"))
	}
	return port
}

func new() *inserter {
	return &inserter{
		db:  sqlx.MustConnect("postgres", *dataSourceName),
		ser: getPort(),
	}
}

//poll forever will look for newlines and assume whatever it got was a full on SQL query worthy of executing.
func (ins *inserter) poll() {
	ins.db.MustExec(fmt.Sprintf(create, *table))
	rdr := bufio.NewReader(ins.ser)
	lines := 0
	for {
		line, err := rdr.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		lines++
		fmt.Printf("%d lines ingested\r", lines)
		if strings.Contains(line, "INSERT INTO") {
			go ins.insert(line)
		}
	}
}

func (ins *inserter) insert(data string) {
	data = fmt.Sprintf(data, *table)
	if _, err := ins.db.Exec(data); err != nil {
		fmt.Println(data)
		fmt.Println(errors.Wrap(err, "Unable to insert"))
		return
	}
}

func main() {
	app.Parse(os.Args[1:])
	i := new()
	i.poll()
}
