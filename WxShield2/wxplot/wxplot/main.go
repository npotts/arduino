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
	"github.com/alecthomas/kingpin"
	"github.com/npotts/arduino/WxShield2/wxplot"
	"os"
)

var (
	app            = kingpin.New("wxPlot", "Plot data pulled from postgres/cockroachdb")
	dataSourceName = app.Flag("dataSource", "Where should we connect to and yelp at (usually a string like 'postgres://user:password@server/database?sslmode=disable')").Short('s').Default("postgres://wx:wx@pika/wx?sslmode=disable").String()
	table          = app.Flag("table", "The database table read raw data from").Short('t').Default("raw").String()
	dir            = app.Flag("output-dir", "Where should the output files be shoveled").Short('p').Default(".").String()
)

func main() {
	if _, err := app.Parse(os.Args[1:]); err != nil {
		panic(err)
	}
	i := wxplot.New(*dataSourceName, *table)
	i.WriteFile(*dir+"/hourly.html", i.Hourly())
	i.WriteFile(*dir+"/weekly.html", i.Weekly())
	i.WriteFile(*dir+"/monthly.html", i.Monthly())
}
