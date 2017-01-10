package plots

import (
	"github.com/gonum/plot/plotter"
	"image/color"
)

/*Trace is a single trace's data*/
type Trace struct {
	Data  plotter.XYs
	Color color.RGBA
}
