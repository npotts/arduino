package wxplot

import (
	"bytes"
	"database/sql"
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

type dataset struct {
	Variable string
	Frames   frames
}

var datasetTmpl = `{{.Variable}}: {
	datasets: [
		// {label: "Pressure", data: [
		// 	{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.Pressure}} },{{end}}
		// ]},
		{label: "Temp A", data: [
			{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.Tempa}} },{{end}}
		]},
		{label: "Temp B", data: [
			{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.Tempb}} },{{end}}
		]},
		{label: "Humidity", data: [
			{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.Humidity}} },{{end}}
		]},
		{label: "PTemp", data: [
			{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.PTemp}} },{{end}}
		]},
		{label: "HTemp", data: [
			{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.HTemp}} },{{end}}
		]},
		{label: "Humidity", data: [
			{{range .Frames}}{x: moment("{{.Timestamp.Format "2006-01-02 15:04:05.999999"}}"), y: {{.Tempb}} },{{end}}
		]}
	]
}`

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
	t := template.Must(template.New("generic").Parse(datasetTmpl))
	ds := dataset{Variable: variable, Frames: *f}
	buf := &bytes.Buffer{}
	t.Execute(buf, ds)
	return buf.String()
}
