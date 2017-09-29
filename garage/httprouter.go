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
	"context"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"time"
)

/*Daemon monitors the garage door, and stabs mice that run around*/
type Daemon struct {
	ctx           context.Context
	cncl          context.CancelFunc
	p             *Parser
	mux           *mux.Router
	box           *rice.Box
	svr           *http.Server
	after, before time.Time
}

func (d *Daemon) html(w http.ResponseWriter, r *http.Request) {
	b := d.box.MustBytes("index.html")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func shouldCloseDoor(a, z time.Time) bool {
	now := time.Now()
	c := time.Date(0, 0, 0, now.Hour(), now.Minute(), now.Second(), 0, time.Local)
	return c.After(a) || c.Before(z)
}

func (d *Daemon) run() {
	ctx, cancel := context.WithCancel(context.Background())
	d.svr.ConnState = func(con net.Conn, st http.ConnState) {
		switch st {
		case http.StateIdle:
		case http.StateClosed:
			cancel()
		}
	}
	go func() {
		for {
			select {
			case <-time.After(1 * time.Minute):
				st := d.p.Next()
				if shouldCloseDoor(d.after, d.before) && !st.Closed {
					fmt.Println("Automatically closing door")
					d.p.issue(closeDoor)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

/*Run starts the server*/
func (d *Daemon) Run() error {
	d.run()
	err := d.svr.ListenAndServe()
	d.cncl() //cancel ctx chain
	fmt.Printf("Server Error: %v\n", err)
	return err
}

/*NewDaemon returns an intialized daemon instance, or panics*/
func NewDaemon(closeAfter, closeBefore, device string, baud, webport int) *Daemon {
	ctx, cancel := context.WithCancel(context.Background())

	panicif := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	after, err0 := time.Parse(time.Kitchen, closeAfter)
	panicif(err0)
	before, err1 := time.Parse(time.Kitchen, closeBefore)
	panicif(err1)

	//open Parser
	port := fmt.Sprintf("serial://%s:%d", device, baud)
	p, err2 := NewParser(ctx, port)
	panicif(err2)
	fmt.Println(port, "opened")

	// Waiting because when opening arduino, it gets reset and take
	// some time to wake
	<-time.After(2 * time.Second)
	defer p.Start()

	m := mux.NewRouter()

	d := &Daemon{
		cncl: cancel,
		p:    p,
		box:  rice.MustFindBox("html"),
		mux:  m,
		svr: &http.Server{
			Handler: m,
			Addr:    fmt.Sprintf(":%d", webport),
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
		after:  after,
		before: before,
	}

	api := m.Path("/api").Subrouter()
	api.Methods("GET").HandlerFunc(d.p.Status)
	api.Methods("PUT").HandlerFunc(d.p.Button)
	api.Methods("POST").HandlerFunc(d.p.OpenDoor)
	api.Methods("DELETE").HandlerFunc(d.p.CloseDoor)
	api.Methods("PATCH").HandlerFunc(d.p.Recal)

	m.HandleFunc("/{*.html}", d.html).Methods("GET")
	m.HandleFunc("/{*.htm}", d.html).Methods("GET")
	m.HandleFunc("/", d.html).Methods("GET")

	return d
}
