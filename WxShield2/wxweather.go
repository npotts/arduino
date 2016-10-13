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

package wxshield2

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

/*Packet is a single sample made in time*/
type Packet struct {
	ID        uint64          `sql:"id" json:"id"`
	Timestamp time.Time       `sql:"timestamp" json:"timestamp"`
	Pressure  sql.NullFloat64 `sql:"pressure" json:"pressure"`
	Tempa     sql.NullFloat64 `sql:"tempa" json:"tempa"`
	Tempb     sql.NullFloat64 `sql:"tempb" json:"tempb"`
	Humidity  sql.NullFloat64 `sql:"humidity" json:"humidity"`
	PTemp     sql.NullFloat64 `sql:"ptemp" json:"ptemp"`
	HTemp     sql.NullFloat64 `sql:"htemp" json:"htemp"`
	Battery   sql.NullFloat64 `sql:"battery" json:"battery"`
	Indx      int             `sql:"indx" json:"indx"`
}

/*tags returns the tags used in the packet structure*/
func (p Packet) tags() []string {
	r := []string{}
	s := reflect.ValueOf(&p).Elem()
	typeof := s.Type()
	for i := 0; i < s.NumField(); i++ {
		r = append(r, typeof.Field(i).Tag.Get("sql"))
	}
	return r
}

/*InsertEroteme uses the sqlite/postgres ? placeholder*/
func (p Packet) InsertEroteme(table string) string {
	tags := p.tags()
	return fmt.Sprintf("INSERT (%s) INTO %s (?%s)", strings.Join(tags, ", "), table, strings.Repeat(", ?", len(tags)-1))
}

/*InsertNamed returns a named SQL query to insert named values*/
func (p Packet) InsertNamed(table string) string {
	tags, named := p.tags(), []string{}
	for _, tag := range tags {
		named = append(named, fmt.Sprintf(":%s", tag))
	}
	return fmt.Sprintf("INSERT (%s) INTO %s (?%s)", strings.Join(tags, ", "), table, strings.Join(named, ", "))
}
