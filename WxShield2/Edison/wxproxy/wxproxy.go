package main

import (
	"bufio"
	"fmt"
	"github.com/alecthomas/kingpin"
	_ "github.com/cockroachdb/pq"
	"github.com/dustin/go-rs232"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
)

var (
	create = "CREATE TABLE %s.%s (id SERIAL, timestamp TIMESTAMP DEFAULT now(), pressure FLOAT, tempa FLOAT, tempb FLOAT, humidity FLOAT, ptemp FLOAT, htemp FLOAT, battery FLOAT, indx INT)"

	app            = kingpin.New("wxProxy", "Shovel data from an arduino into a cockroachdb somewhere else")
	dataSourceName = app.Flag("dataSource", "Where should we connect to and yelp at (usually a string like 'postgresql://root@dataserver:26257?sslmode=disable')").Short('s').Default("postgresql://root@chipmunk:26257?sslmode=disable").String()
	database       = app.Flag("database", "The database to aim at").Short('d').Default("wx").String()
	raw            = app.Flag("table", "The database table to fire into").Short('t').Default("raw").String()
	device         = app.Flag("device", "The RS232 serial device to read from (at 57600 baud)").Short('D').Default("/dev/ttyUSB0").String()
)

type inserter struct {
	db  *sqlx.DB
	ser *rs232.SerialPort
}

func mustOpen(port *rs232.SerialPort, err error) *rs232.SerialPort {
	if err != nil {
		panic(errors.Wrap(err, "Unable to open serial device"))
	}
	return port
}

func new() *inserter {
	i := &inserter{
		db:  sqlx.MustConnect("postgres", *dataSourceName),
		ser: mustOpen(rs232.OpenPort(*device, 57600, rs232.S_8N1)),
	}
	i.db.MustExec(fmt.Sprintf(create, *database, *raw))
	return i
}

//poll forever will look for newlines and assume whatever it got was a full on SQL query worthy of executing.
func (ins *inserter) poll() {
	rdr := bufio.NewReader(ins.ser)
	for {
		line, err := rdr.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		go ins.insert(line)
	}
}

func (ins *inserter) insert(data string) {
	data = strings.Replace(data, "wx", fmt.Sprintf("%d.%d", *database, *raw), 1)
	if _, err := ins.db.Exec(data); err != nil {
		fmt.Println(errors.Wrap(err, "Unable to insert"))
	}
}

func main() {
	kingpin.Parse()
	i := new()
	i.poll()
}
