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
	"fmt"
	"testing"
)

func TestStatus_ParseString(t *testing.T) {
	s := &Status{}
	if e := s.ParseString(`243264,802,200,700,100,0,1
>>
243352,801,200,700,100,0,1
>>
243438,802,200,700,100,0,1
>>
243522,800,200,700,100,0,1
>>
243609,801,200,700,100,0,1
>>
243697,803,200,700,100,0,1
>>
243782,801,200,700,100,0,1
>>
243867,804,200,700,100,0,1
>>
243953,803,200,700,100,0,1
>>
244038,800,200,700,100,0,1
>>
244123,804,200,700,100,0,1
>>`); e != nil {
		t.Errorf("Error: Got %v", e)
	}

	if e := s.ParseString(`243352,801,200,700,100,0`); e == nil {
		t.Error("SHould have gotten garbage")
	}
	if e := s.ParseString(`243352,801,200,700,100,0,a`); e == nil {
		t.Error("SHould have gotten garbage")
	}
	if e := s.ParseString(`243352,801,200,700,100,a,1`); e == nil {
		t.Error("SHould have gotten garbage")
	}
	if e := s.ParseString(`243352,801,200,700,a,0,1`); e == nil {
		t.Error("SHould have gotten garbage")
	}
	if e := s.ParseString(`243352,801,200,a,100,0,1`); e == nil {
		t.Error("SHould have gotten garbage")
	}
	if e := s.ParseString(`243352,801,a,700,100,0,1`); e == nil {
		t.Error("SHould have gotten garbage")
	}
	if e := s.ParseString(`243352,a,200,700,100,0,1`); e == nil {
		t.Error("SHould have gotten garbage")
	}
	if e := s.ParseString(`a,801,200,700,100,0,1`); e == nil {
		t.Error("SHould have gotten garbage")
	}
}

func TestStatus_Parse(t *testing.T) {
	s := &Status{}
	if e := s.Parse([]string{"35103", "612", "200", "700", "82", "0", "0"}); e != nil {
		t.Fatalf("Didnt parse as expected %v", e)
	}
	if e := s.Parse([]string{"a"}); e == nil {
		t.Error("Should have goten error for invalid length")
	}
}

func TestNewParser(t *testing.T) {
	if !testing.Short() {
		t.Skip("Not cehcking hardware")
		t.SkipNow()
	}
	s, err := NewParser("serial:///dev/cu.usbserial-A6004nZq:115200")
	if err != nil {
		t.Errorf("Unable to dial: %v", err)
		t.FailNow()
	}

	for i := 0; i < 40; i++ {
		fmt.Println(s.Next())
	}
}
