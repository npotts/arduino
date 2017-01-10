package plots

import (
	"io"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

/*NewTimePlot returns a partially initiated TimePlot*/
func NewTimePlot(title, xlabel, ylabel string) *TimePlot {
	return &TimePlot{
		title:  title,
		xlabel: xlabel,
		ylabel: ylabel,
		plots:  []Trace{},
	}
}

/*TimePlot ...*/
type TimePlot struct {
	title, xlabel, ylabel string
	plots                 []Trace
}

/*AddTrace adds a trace to time plots*/
func (tp *TimePlot) AddTrace(trace Trace) {
	tp.plots = append(tp.plots, trace)
}

/*WriteTo writes to the writer the the plots configued with a width and height with a output format*/
func (tp *TimePlot) WriteTo(w io.Writer, width vg.Length, height vg.Length, format string) error {
	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = tp.title
	p.X.Label.Text = tp.xlabel
	p.Y.Label.Text = tp.ylabel
	p.X.Tick.Marker = plot.TimeTicks{Format: "Jan 02\n15:04"}
	for _, trace := range tp.plots {
		line, points, err := plotter.NewLinePoints(trace.Data)
		if err != nil {
			return err
		}
		points.Shape = draw.CircleGlyph{}
		line.Color = trace.Color
		points.Color = trace.Color
		p.Add(line, points)
	}

	wto, err := p.WriterTo(width, height, format)
	if err != nil {
		return err
	}
	_, err = wto.WriteTo(w)
	return err
}
