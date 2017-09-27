/*
The MIT License (MIT)

Copyright (c) 2016-2017 Nick Potts

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

package garage

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/NCAR/ACSd/monitor/arbiter"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var (
	alwaysSuccess = arbiter.Command{
		Name:          "Always succeeds",
		Timeout:       1000 * time.Millisecond,
		Prototype:     "",
		CommandRegexp: regexp.MustCompile(".*"),
		Response:      regexp.MustCompile(".*"),
		Error:         nil,
		Description:   "Always succeeds",
	}
	cstatus = arbiter.Command{
		Name:          "Status",
		Timeout:       500 * time.Millisecond,
		Prototype:     "?",
		CommandRegexp: regexp.MustCompile("\\?"),
		Response:      regexp.MustCompile("[\\d]*,[\\d]*,[\\d]*,[\\d]*,[\\d]*,[\\d]*,[\\d]*\r\n"),
		Error:         nil,
		Description:   "Get status",
	}
	recal = arbiter.Command{
		Name:          "Recal",
		Timeout:       2 * time.Minute,
		Prototype:     "~",
		CommandRegexp: regexp.MustCompile("~"),
		Response:      regexp.MustCompile("[\\d]*,[\\d]*,[\\d]*,[\\d]*,[\\d]*,[\\d]*,[\\d]*\r\n"),
		Error:         nil,
		Description:   "Get status",
	}
	open = arbiter.Command{
		Name:          "Open Door",
		Timeout:       20 * time.Second,
		Prototype:     "o",
		CommandRegexp: regexp.MustCompile("o"),
		Response:      regexp.MustCompile("opened\n"),
		Error:         regexp.MustCompile("nope\n"),
		Description:   "Opens the Door",
	}
	close = arbiter.Command{
		Name:          "Close Door",
		Timeout:       20 * time.Second,
		Prototype:     "c",
		CommandRegexp: regexp.MustCompile("c"),
		Response:      regexp.MustCompile("closed\n"),
		Error:         regexp.MustCompile("nope\n"),
		Description:   "Closes the door",
	}
	toggle = arbiter.Command{
		Name:          "Toggle Button",
		Timeout:       10 * time.Second,
		Prototype:     "^",
		CommandRegexp: regexp.MustCompile("^"),
		Response:      regexp.MustCompile("Pushing button\n"),
		Error:         nil,
		Description:   "Toggles the door",
	}
)

/*Status is the message the garage door emits*/
type Status struct {
	Millis      uint64 `json:"millis"`       //64 bits
	Pos         uint64 `json:"pos"`          //16 bits
	FloorADC    uint64 `json:"floor_adc"`    //16bits
	CeilingADC  uint64 `json:"ceiling_adc"`  //16bits
	PercentOpen uint64 `json:"percent_open"` //8bits
	Closed      bool   `json:"closed"`
	FullyOpen   bool   `json:"fully_open"`
}

/*ParseRaw parses some raw bytes for a valid status message*/
func (s *Status) ParseRaw(raw []byte) error {
	csvp := csv.NewReader(bytes.NewBuffer(raw))
	csvp.Comma = ','
	csvp.Comment = '>'
	csvp.FieldsPerRecord = 7
	csvp.LazyQuotes = false
	csvp.TrimLeadingSpace = true
	recs, err := csvp.ReadAll()
	if err != nil {
		return err
	}
	for _, rec := range recs {
		err = s.Parse(rec)
	}
	return err
}

/*ParseString parses the passed string for status bits*/
func (s *Status) ParseString(parseString string) error {
	return s.ParseRaw([]byte(parseString))
}

/*Parse takes a split list of string slices for values*/
func (s *Status) Parse(parsed []string) (err error) {
	if len(parsed) != 7 {
		return errors.New("Invalid number of parameters")
	}
	if s.Millis, err = strconv.ParseUint(parsed[0], 10, 64); err != nil {
		return err
	}
	if s.Pos, err = strconv.ParseUint(parsed[1], 10, 16); err != nil {
		return err
	}
	if s.FloorADC, err = strconv.ParseUint(parsed[2], 10, 16); err != nil {
		return err
	}
	if s.CeilingADC, err = strconv.ParseUint(parsed[3], 10, 16); err != nil {
		return err
	}
	if s.PercentOpen, err = strconv.ParseUint(parsed[4], 10, 8); err != nil {
		return err
	}
	if s.Closed, err = strconv.ParseBool(parsed[5]); err != nil {
		return err
	}
	if s.FullyOpen, err = strconv.ParseBool(parsed[6]); err != nil {
		return err
	}
	return nil
}

/*String conforms to the Stringer interface*/
func (s Status) String() string {
	return fmt.Sprintf(`Clk: %d P:% 4d Floor:% 4d Ceiling:% 4d %%Open:% 3d Closed: %5v FullyOpen: %5v`, s.Millis, s.Pos, s.FloorADC, s.CeilingADC, s.PercentOpen, s.Closed, s.FullyOpen)
}

/*NewParser returns a new parsers after openening the arbiter*/
func NewParser(dial string) (*Parser, error) {
	arb, err := arbiter.OpenArbiter(dial, 1000*time.Millisecond, alwaysSuccess)
	return &Parser{dev: arb}, err
}

/*Parser wraps around the arduino and provides a simple way to get status
and perform limited control on a garage door opener*/
type Parser struct {
	mux sync.Mutex
	dev arbiter.Arbiter
}

func (p *Parser) issue(cmd arbiter.Command) ([]byte, error) {
	p.mux.Lock()
	defer p.mux.Unlock()
	resp := p.dev.Control(cmd)
	return resp.Bytes, resp.Error
}

/*Next polls the device for another status message*/
func (p *Parser) Next() (rtn *Status, err error) {
	b, e := p.issue(cstatus)
	if e == nil {
		rtn = &Status{}
		return rtn, rtn.ParseRaw(cstatus.Response.FindSubmatch(b)[0])
	}
	return nil, e
}

/*Status polls for the next status message and returns the JSON equivalent of it*/
func (p *Parser) Status(w http.ResponseWriter, r *http.Request) {
	s, e := p.Next()
	if e == nil {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(s)
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%s", e)
}

/*OpenDoor Opens the doore*/
func (p *Parser) OpenDoor(w http.ResponseWriter, r *http.Request) {
	go func() {
		_, e := p.issue(open)
		if e != nil {
			fmt.Println("Error Opening Door: ", e)
		}
	}()
	w.WriteHeader(http.StatusNoContent)
}

/*CloseDoor closes the doore*/
func (p *Parser) CloseDoor(w http.ResponseWriter, r *http.Request) {
	go func() {
		_, e := p.issue(close)
		if e != nil {
			fmt.Println("Error closing Door: ", e)
		}
	}()
	w.WriteHeader(http.StatusNoContent)
}

/*Recal issues a recal*/
func (p *Parser) Recal(w http.ResponseWriter, r *http.Request) {
	go func() {
		_, e := p.issue(recal)
		if e != nil {
			fmt.Println("Error issuing recal: ", e)
		}
	}()
	w.WriteHeader(http.StatusNoContent)
}

/*Button punches the button*/
func (p *Parser) Button(w http.ResponseWriter, r *http.Request) {
	go func() {
		_, e := p.issue(toggle)
		if e != nil {
			fmt.Println("Error issuing toggle: ", e)
		}
	}()
	w.WriteHeader(http.StatusNoContent)
}
