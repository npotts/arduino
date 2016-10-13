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
	ID        uint64          `sql:"-"`
	Timestamp time.Time       `sql:"timestamp"`
	Pressure  sql.NullFloat64 `sql:"pressure"`
	Tempa     sql.NullFloat64 `sql:"tempa"`
	Tempb     sql.NullFloat64 `sql:"tempb"`
	Humidity  sql.NullFloat64 `sql:"humidity"`
	PTemp     sql.NullFloat64 `sql:"ptemp"`
	HTemp     sql.NullFloat64 `sql:"htemp"`
	Battery   sql.NullFloat64 `sql:"battery"`
	Indx      int             `sql:"indx"`
}

/*jsonPacket is a json encodable packet*/
type jsonPacket struct {
	Pressure *float64 `json:"pressure"`
	Tempa    *float64 `json:"tempa"`
	Tempb    *float64 `json:"tempb"`
	Humidity *float64 `json:"humidity"`
	PTemp    *float64 `json:"ptemp"`
	HTemp    *float64 `json:"htemp"`
	Battery  *float64 `json:"battery"`
	Indx     *int     `json:"indx"`
}

/*tags returns the tags used in the packet structure*/
func (p Packet) tags() []string {
	r := []string{}
	s := reflect.ValueOf(&p).Elem()
	typeof := s.Type()
	for i := 0; i < s.NumField(); i++ {
		if t := typeof.Field(i).Tag.Get("sql"); t != "-" {
			r = append(r, t)
		}
	}
	return r
}

/*Jsonable returns a structure that can be written to via json*/
func (p *Packet) Jsonable() *jsonPacket {
	p.Pressure.Valid, p.Tempa.Valid, p.Tempb.Valid, p.Humidity.Valid, p.PTemp.Valid, p.HTemp.Valid, p.Battery.Valid = true, true, true, true, true, true, true
	return &jsonPacket{
		Pressure: &p.Pressure.Float64,
		Tempa:    &p.Tempa.Float64,
		Tempb:    &p.Tempb.Float64,
		Humidity: &p.Humidity.Float64,
		PTemp:    &p.PTemp.Float64,
		HTemp:    &p.HTemp.Float64,
		Battery:  &p.Battery.Float64,
		Indx:     &p.Indx,
	}
}

/*InsertEroteme uses the sqlite/postgres ? placeholder*/
func (p Packet) InsertEroteme(table string) string {
	tags := p.tags()
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (?%s)", table, strings.Join(tags, ", "), strings.Repeat(", ?", len(tags)-1))
}

/*InsertNamed returns a named SQL query to insert named values*/
func (p Packet) InsertNamed(table string) string {
	tags, named := p.tags(), []string{}
	for _, tag := range tags {
		named = append(named, fmt.Sprintf(":%s", tag))
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(tags, ", "), strings.Join(named, ", "))
}
