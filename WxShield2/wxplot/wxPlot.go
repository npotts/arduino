package wxplot

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"

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
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE timestamp > $1 AND timestamp < $2", p.database, p.table)
	err = p.db.Select(&f, query, start, end)
	return
}

/*Hourly geerates 24-hour plots that start at midnight the day before now and end at midnight*/
func (p PlotUtil) Hourly() {
	end := time.Now()
	start := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.Local)
	data, _ := p.databyRange(start, end)
	fmt.Println(data.html("Daily"))
}

/*Weekly generates 7-day plots*/
func (p PlotUtil) Weekly() {
	end := time.Now()
	start := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, nil)
	data, _ := p.databyRange(start, end)
	fmt.Println(data.html("Daily"))
}
