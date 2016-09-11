package main

import (
	"bufio"
	"fmt"
	"github.com/alecthomas/kingpin"
	_ "github.com/cockroachdb/pq"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/tarm/serial"
	"os"
	"strings"
)

var (
	create = "CREATE TABLE IF NOT EXISTS %s.%s  (id SERIAL, timestamp TIMESTAMP DEFAULT now(), pressure FLOAT, tempa FLOAT, tempb FLOAT, humidity FLOAT, ptemp FLOAT, htemp FLOAT, battery FLOAT, indx INT)"

	app            = kingpin.New("wxProxy", "Shovel data from an arduino into a postgres/cockroachdb")
	dataSourceName = app.Flag("dataSource", "Where should we connect to and yelp at (usually a string like 'postgresql://root@dataserver:26257?sslmode=disable')").Short('s').Default("postgresql://root@chipmunk:26257?sslmode=disable").String()
	database       = app.Flag("database", "The database to aim at").Short('d').Default("wx").String()
	raw            = app.Flag("table", "The database table to fire into").Short('t').Default("raw").String()
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
	ins.db.MustExec(fmt.Sprintf(create, *database, *raw))
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
	data = fmt.Sprintf(data, *database, *raw)
	if _, err := ins.db.Exec(data); err != nil {
		fmt.Println(errors.Wrap(err, "Unable to insert"))
		return
	}
}

func main() {
	app.Parse(os.Args[1:])
	i := new()
	i.poll()
}
