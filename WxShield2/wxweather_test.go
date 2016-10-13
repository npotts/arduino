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
	"testing"
)

func TestPacket_InsertEroteme(t *testing.T) {
	n := Packet{}
	want := `INSERT (id, timestamp, pressure, tempa, tempb, humidity, ptemp, htemp, battery, indx) INTO test (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if n.InsertEroteme("test") != want {
		t.Errorf("Got %q not %q", n.InsertEroteme("test"), want)
	}
}

func TestPacket_InsertNamed(t *testing.T) {
	n := Packet{}
	want := `INSERT (id, timestamp, pressure, tempa, tempb, humidity, ptemp, htemp, battery, indx) INTO test (?:id, :timestamp, :pressure, :tempa, :tempb, :humidity, :ptemp, :htemp, :battery, :indx)`
	if n.InsertNamed("test") != want {
		t.Errorf("Got %q not %q", n.InsertNamed("test"), want)
	}
}
