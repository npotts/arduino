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
	//find out how many samples we are dealing with
	i, samples := 0, 750
	if err = p.db.Select(&i, fmt.Sprintf("SELECT COUNT(*) FROM %s.%s", p.database, p.table)); err != nil {
		return
	}
	i = int(i / samples) //max of 750ish samples
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE timestamp > $1 AND timestamp < $2  ORDER BY timestamp LIMIT %d OFFSET %d;", p.database, p.table, samples, i)
	err = p.db.Select(&f, query, start, end)
	return
}

//Hourly geerates 24-hour
func (p PlotUtil) Hourly() {
	end := time.Now()
	start := end.Add(-1 * time.Hour)
	data, err := p.databyRange(start, end)
	fmt.Println(err)
	fmt.Println(len(data))
	fmt.Println(data[0])
	fmt.Println(data[len(data)-1:])
	fmt.Println(data.plotworthy("yippee"))
}
