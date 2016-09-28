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

package wxplot

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"time"

	_ "github.com/cockroachdb/pq" //Postgres sql hook
)

//PlotUtil generated HTML plots from data stred on database
type PlotUtil struct {
	db              *sqlx.DB
	database, table string
}

/*New returns a created PlotUtil or panics trying*/
func New(dataSourceName, table string) *PlotUtil {
	return &PlotUtil{
		db:    sqlx.MustConnect("postgres", dataSourceName),
		table: table,
	}
}

func (p *PlotUtil) databyRange(start, end time.Time) (f frames, err error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE timestamp > $1 AND timestamp < $2  ORDER BY timestamp", p.table)
	if err = p.db.Select(&f, query, start, end); err != nil {
		return frames{}, err
	}
	//probably not the most efficient manner here
	decimate := len(f) / 250
	newf := frames{}
	for i, frame := range f {
		if i%decimate == 0 {
			frame.Timestamp = frame.Timestamp.Local()
			newf = append(newf, frame)
		}
	}
	return newf, nil
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
	start = start.Add(-24 * time.Duration(end.Weekday()) * time.Hour)
	data, _ := p.databyRange(start.UTC(), end.UTC())
	return data.html("Weekly")
}

/*Monthly generates Calandar month plots*/
func (p PlotUtil) Monthly() string {
	end := time.Now().UTC().Local()
	start := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, end.Location())
	data, _ := p.databyRange(start.UTC(), end.UTC())
	return data.html("Monthy")
}

func (p *PlotUtil) WriteFile(filename, data string) {
	ioutil.WriteFile(filename, []byte(data), 0777)
}
