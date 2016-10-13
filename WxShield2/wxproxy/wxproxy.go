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
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/tarm/serial"
	"os"

	_ "github.com/lib/pq"
	"github.com/npotts/arduino/WxShield2"
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
	fmt.Printf("Attempting to open %s: ", *device)
	port, err := serial.OpenPort(&serial.Config{Name: *device, Baud: 57600, Size: 8, Parity: serial.ParityNone, StopBits: 1})
	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to open serial device"))
		os.Exit(1)
	}
	fmt.Println("done")
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
	fmt.Println("Entering Polling Loop")
	for {
		line, err := rdr.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		packet := &wxshield2.Packet{}
		err = json.Unmarshal(line, packet.Jsonable())
		if err == nil {
			go ins.insert(*packet)
			lines++
			fmt.Printf("\r%d lines ingested", lines)
		}
	}
}

/*insert exec*/
func (ins *inserter) insert(packet wxshield2.Packet) {
	sql := packet.InsertNamed(*table)
	if _, err := ins.db.NamedExec(sql, packet); err != nil {
		fmt.Println(errors.Wrapf(err, "Unable to insert: %q", sql))
		return
	}
}

func main() {
	app.Parse(os.Args[1:])
	i := new()
	i.poll()
}
