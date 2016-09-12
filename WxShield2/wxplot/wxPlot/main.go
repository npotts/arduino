package main

import (
	"github.com/alecthomas/kingpin"
	_ "github.com/cockroachdb/pq"
	"github.com/npotts/arduino/WxShield2/wxplot"
	"os"
)

var (
	app            = kingpin.New("wxPlot", "Plot data pulled from postgres/cockroachdb")
	dataSourceName = app.Flag("dataSource", "Where should we connect to and yelp at (usually a string like 'postgresql://root@dataserver:26257?sslmode=disable')").Short('s').Default("postgresql://root@chipmunk:26257?sslmode=disable").String()
	database       = app.Flag("database", "The database to aim at").Short('d').Default("wx").String()
	raw            = app.Flag("table", "The database table read raw data from").Short('t').Default("raw").String()
	// daily          = app.Flag("daily", "").Short('a').Default("true").Bool()
)

func main() {
	app.Parse(os.Args[1:])
	i := wxplot.New(dataSourceName, database, raw)
	i.Hourly()
	// i.poll()
}
