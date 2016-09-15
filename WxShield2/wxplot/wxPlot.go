package wxplot

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
  "io/ioutil"

	_ "github.com/cockroachdb/pq" //Postgres sql hook
)

//PlotUtil generated HTML plots from data stred on database
type PlotUtil struct {
	db              *sqlx.DB
	database, table string
}

/*New returns a created PlotUtil or panics trying*/
func New(dataSourceName, database, table string) *PlotUtil {
	return &PlotUtil{
		db:       sqlx.MustConnect("postgres", dataSourceName),
		database: database,
		table:    table,
	}
}

func (p *PlotUtil) databyRange(start, end time.Time) (f frames, err error) {
	//find out how many samples we are dealing with
	i, samples := 0, 750
	if err = p.db.Get(&i, fmt.Sprintf("SELECT COUNT(*) FROM %s.%s WHERE timestamp > $1 and timestamp < $2", p.database, p.table), start, end); err != nil {
		return
	}
 	i = int(i / samples) //max of 750ish samples
	//query := fmt.Sprintf("SELECT * FROM %s.%s WHERE timestamp > $1 AND timestamp < $2  ORDER BY timestamp LIMIT %d OFFSET %d;", p.database, p.table, samples, i)
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE timestamp > $1 AND timestamp < $2  ORDER BY timestamp", p.database, p.table)
	err = p.db.Select(&f, query, start, end)
  return
}

/*Hourly geerates 24-hour plots that start at midnight the day before now and end at midnight*/
func (p PlotUtil) Hourly() string {
	end := time.Now()
	start := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.Local).UTC()
	data, _ := p.databyRange(start.UTC(), end)
  return data.html("Daily")
}

/*Weekly generates 7-day plots*/
func (p PlotUtil) Weekly() string {
  end := time.Now().UTC().Local()
  start := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
  start = start.Add(-24 * time.Duration(end.Weekday()) * time.Hour);
	data, _ := p.databyRange(start.UTC(), end.UTC())
	return data.html("Weekly")
}

func (p *PlotUtil) WriteFile(filename, data string) {
  ioutil.WriteFile(filename, []byte(data), 0777);
}
