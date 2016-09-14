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

//Hourly geerates 24-hour
func (p PlotUtil) Hourly() {
	end := time.Now()
	start := end.Add(-2 * time.Minute)
	data, _ := p.databyRange(start, end)
	fmt.Println(data.html("Wala"))
	// fmt.Println(err)
	// fmt.Println(len(data))
	// fmt.Println(data[0])
	// fmt.Println(data[len(data)-1:])
	// fmt.Println(data.plotworthy("yippee"))
}
