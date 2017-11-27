// +build ignore

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
	"github.com/alecthomas/kingpin"
	"github.com/npotts/arduino/garage"
	"os"
)

var (
	app     = kingpin.New("garage", "Exporting a garage door to the great wide internets")
	baud    = app.Flag("baud", "Specify the baud rate. Default is the compiled in baud rate").Short('b').Default("115200").Int()
	device  = app.Flag("device", "The RS232 serial device connected to the Arduino that controls the garage door (http://github.com/npotts/arduino/garage)").Short('d').Default("serial:///dev/ttyUSB0").String()
	webport = app.Flag("port", "HTTP port to listen for incoming connections from").Short('p').Default("80").Int()
	after   = app.Flag("after", "If the garage is open after this time, in Kitchen Timer format (11:32AM), the daemon will attempt to close the door automatically.").Short('t').Default("10:00PM").String()
	before  = app.Flag("before", "Same as above, but before this.  Sorry if you work the graveyard shift").Short('T').Default("7:00AM").String()
)

func main() {
	app.Parse(os.Args[1:])
	garage.NewDaemon(*after, *before, *device, *baud, *webport).Run()
}
