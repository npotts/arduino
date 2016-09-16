package main

import (
	"github.com/alecthomas/kingpin"
	"github.com/npotts/arduino/WxShield2/wxplot"
	"os"
)

var (
	app            = kingpin.New("wxPlot", "Plot data pulled from postgres/cockroachdb")
	dataSourceName = app.Flag("dataSource", "Where should we connect to and yank data from (usually a string like 'postgresql://root@dataserver:26257?sslmode=disable')").Short('s').Default("postgresql://root@chipmunk:26257?sslmode=disable").String()
	database       = app.Flag("database", "The database to aim at").Short('d').Default("wx").String()
	raw            = app.Flag("table", "The database table read raw data from").Short('t').Default("raw").String()
	dir            = app.Flag("output-dir", "Where should the output files be shoveled").Short('p').Default(".").String()
)

func main() {
	if _, err := app.Parse(os.Args[1:]); err != nil {
		panic(err)
	}
	i := wxplot.New(*dataSourceName, *database, *raw)
	i.WriteFile(*dir+"/hourly.html", i.Hourly())
	i.WriteFile(*dir+"/weekly.html", i.Weekly())
}
