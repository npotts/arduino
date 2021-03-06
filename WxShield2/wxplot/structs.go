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
	"bytes"
	"database/sql"
	"fmt"
	"text/template"
	"time"
)

type frame struct {
	ID        int             `sql:"id"`
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

type frames []frame

type singlePlot struct {
	Fieldname string
	Label     string
	Solid     string
	Opaque    string
	Frames    frames
}

var singlePlotTmpl = `{
                        label: "{{.Label}}",
                        fill: false,
                        borderColor: "{{.Solid}}",
                        backgroundColor: "{{.Opaque}}",
                        pointBorderColor: "{{.Opaque}}",
                        pointBackgroundColor: "{{.Opaque}}",
                        pointBorderWidth: 1,
                        data: [{{range .Frames}}
                            {x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.%s.Float64}}},{{end}}
                        ]
                    }`

type plotConfig struct {
	Varname string
	Title   string
	YUnits  string
	Data    []string
}

var plotConfigTmpl = `var {{.Varname}} = {
            type: "line",
            options: {
                responsive: true,
                title:{
                    display:true,
                    text:"{{.Title}}"
                },
                scales: {
                    xAxes: [{
                        type: "time",
                        display: true,
                        scaleLabel: {
                            display: true,
                            labelString: "Sample Time"
                        }
                    }],
                    yAxes: [{
                        display: true,
                        scaleLabel: {
                            display: true,
                            labelString: "{{.YUnits}}"
                        }
                    }]
                }
            },
            data: {
                datasets: [
                    {{range .Data }}{{.}},{{end}}
                ]
            }
        };
`

/*getRendered turnes a bunch of templates and data fields specified by a singleplot into an array of properly encoded HTML*/
func (f *frames) getRendered(plots []singlePlot) (r []string) {
	for _, plot := range plots {
		templ := template.Must(template.New(plot.Label).Parse(fmt.Sprintf(singlePlotTmpl, plot.Fieldname)))
		buf := &bytes.Buffer{}
		templ.Execute(buf, plot)
		r = append(r, buf.String())
	}
	return
}

func (f *frames) stringizeThis(variable, title, yunits string, plots []singlePlot) string {
	bigplots := plotConfig{
		Varname: variable,
		Title:   title,
		YUnits:  yunits,
		Data:    f.getRendered(plots),
	}
	tmpl := template.Must(template.New("").Parse(plotConfigTmpl))
	buf := &bytes.Buffer{}
	tmpl.Execute(buf, &bigplots)
	return buf.String()

}

/*Returend a rendered Java script for the Temps across the data*/
func (f *frames) Temps() string {
	plots := []singlePlot{
		singlePlot{Fieldname: "Tempa", Label: "Temperature A", Solid: "rgba(255,0,0,1.0)", Opaque: "rgba(255,0,0,0.8)", Frames: *f},
		singlePlot{Fieldname: "Tempb", Label: "Temperature B", Solid: "rgba(255,50,0,1.0)", Opaque: "rgba(255,50,0,0.8)", Frames: *f},
		singlePlot{Fieldname: "PTemp", Label: "T_pressure", Solid: "rgba(0,0,255,1.0)", Opaque: "rgba(0,0,255,0.8)", Frames: *f},
		singlePlot{Fieldname: "HTemp", Label: "T_humidity", Solid: "rgba(0,255,0,1.0)", Opaque: "rgba(0,255,0,0.8)", Frames: *f},
	}
	return f.stringizeThis("tconfig", "Temperatures", "Degrees (C)", plots)
}

/*Returend a rendered Java script for the Humidity data*/
func (f *frames) Humidity() string {
	return f.stringizeThis("hconfig", "Humidity", "Relative Humidity", []singlePlot{
		singlePlot{
			Fieldname: "Humidity",
			Label:     "Humidity",
			Solid:     "rgba(0, 255, 255,1.0)",
			Opaque:    "rgba(0, 255, 255,0.8)",
			Frames:    *f,
		},
	})
}

/*Returend a rendered js for the pressure data*/
func (f *frames) Pressure() string {
	return f.stringizeThis("pconfig", "Pressure", "mbar", []singlePlot{
		singlePlot{
			Fieldname: "Pressure",
			Label:     "Pressure",
			Solid:     "rgba(2153, 102, 255,1.0)",
			Opaque:    "rgba(153, 102, 255,0.8)",
			Frames:    *f,
		},
	})
}

var htmlTmpl = `
<!doctype html>
<html>

<head>
    <title>Weather Data::{{.Timebase}} </title>
    <script src="http://cdnjs.cloudflare.com/ajax/libs/moment.js/2.13.0/moment.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.2.2/Chart.bundle.min.js"></script>
    <script src="http://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
    <style>
    canvas {
        -moz-user-select: none;
        -webkit-user-select: none;
        -ms-user-select: none;
    }
    </style>
</head>

<body>
    <div style="width:95%;"><canvas id="pcanvas"></canvas></div>
    <div style="width:95%;"><canvas id="tcanvas"></canvas></div>
    <div style="width:95%;"><canvas id="hcanvas"></canvas></div>
    <br>
    <script>
    {{range .Data}}{{.}}
    {{end}}

    window.onload = function() {
        var ctxa = document.getElementById("pcanvas").getContext("2d");
        window.myLine = new Chart(ctxa, pconfig);

        var ctxb = document.getElementById("tcanvas").getContext("2d");
        window.myLine = new Chart(ctxb, tconfig);

        var ctxc = document.getElementById("hcanvas").getContext("2d");
        window.myLine = new Chart(ctxc, hconfig);

    };

</script>
</body>

</html>
`

func (f *frames) html(timebase string) string {
	type h struct {
		Timebase string
		Data     []string
	}

	hh := h{
		Timebase: timebase,
		Data:     []string{f.Pressure(), f.Temps(), f.Humidity()},
	}
	tmpl := template.Must(template.New("html").Parse(htmlTmpl))
	buf := &bytes.Buffer{}
	tmpl.Execute(buf, &hh)
	return buf.String()
}
