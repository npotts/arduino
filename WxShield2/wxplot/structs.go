package wxplot

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
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

var singlePlotTmpl = `{ label: "{{.Label}}", fill: false, borderColor: '{{.Solid}}', backgroundColor: '{{.Opaque}}', pointBorderColor: "{{.Opaque}}", pointBackgroundColor: "{{.Opaque}}", pointBorderWidth: 1, data: [{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.%s.Float64}} },{{end}}]},`

type plotConfig struct {
	Varname string
	Title   string
	YUnits  string
	Data    []string
}

var plotConfigTmpl = `var {{.Varname}} = {
            type: 'line',
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
                            labelString: "{{.YUnits}}""
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

// var datasetTmpl = `{{.Variable}}: {
// 	datasets: [

// 		// {label: "Pressure", data: [
// 		// 	{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.Pressure.Float64}} },{{end}}
// 		// ]},
// 		{label: "Temp A", fill: false,  data: [{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.Tempa.Float64}} },{{end}}]},
// 		{label: "Temp B", fill: false,  data: [{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.Tempb.Float64}} },{{end}}]},
// 		{label: "Humidity", fill: false,  data: [{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.Humidity.Float64}} },{{end}}]},
// 		{label: "PTemp", fill: false,  data: [{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.PTemp.Float64}} },{{end}}]},
// 		{label: "HTemp", fill: false,  data: [{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.HTemp.Float64}} },{{end}}]},
// 		{label: "Humidity", fill: false, data: [{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.Tempb.Float64}} },{{end}}]}
// 	]
// }`

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

/*plotworthy converts data from the native structures to stuff that can be shoved into a file
<var>: {
	datasets: [{
	    label: <label>"Dataset with string point data",
	    data: [{
	        x: moment("2016-09-11 21:26:52.001082"),
	        y: 1.0
	    }, {
	        x: moment("2016-09-11 22:26:49.694226"),
	        y: 2.0
	    }],
	    fill: true
	}]
*/
func (f *frames) plotworthy(variable string) string {
	plots := []singlePlot{
		singlePlot{Fieldname: "Tempa", Label: "Temperature 'A'", Solid: "rgba(255,0,0,1.0)", Opaque: "rgba(255,0,0,0.8)", Frames: *f},
		singlePlot{Fieldname: "Tempb", Label: "Temperature 'B'", Solid: "rgba(255,50,0,1.0)", Opaque: "rgba(255,50,0,0.8)", Frames: *f},
		singlePlot{Fieldname: "PTemp", Label: "T_pressure", Solid: "rgba(0,0,255,1.0)", Opaque: "rgba(0,0,255,0.8)", Frames: *f},
		singlePlot{Fieldname: "HTemp", Label: "T_humidity", Solid: "rgba(0,255,0,1.0)", Opaque: "rgba(0,255,0,0.8)", Frames: *f},
	}

	bigplots := plotConfig{
		Varname: "tconfig",
		Title:   "Temperatures",
		YUnits:  "Degrees (C)",
		Data:    f.getRendered(plots),
	}

	tmpl := template.Must(template.New("").Parse(plotConfigTmpl))
	buf := &bytes.Buffer{}
	fmt.Println(tmpl.Execute(buf, &bigplots))
	return buf.String()
}
