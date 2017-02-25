package plots

import (
	"fmt"
	"github.com/gonum/plot/plotter"
	"github.com/guregu/null"
	"time"
)

/*Measurement is a data measurement*/
type Measurement struct {
	Rowid          null.Float `db:"rowid"`
	Created        null.Time  `db:"created"`
	Battery        null.Float `db:"battery"`
	Humidity       null.Float `db:"humidity"`
	HumidityTemp   null.Float `db:"humidityTemp"`
	Ihumidity      null.Float `db:"ihumidity"`
	IhumidityTemp  null.Float `db:"ihumidityTemp"`
	Pressure       null.Float `db:"pressure"`
	PressureTemp   null.Float `db:"pressureTemp"`
	Temperature    null.Float `db:"temperature"`
	TemperatureExt null.Float `db:"temperatureExt"`
	Vref           null.Float `db:"vref"`
}

/*Measurements is a list of measurements*/
type Measurements []Measurement

/*XYs returns a series of XY points*/
func (ms Measurements) XYs(field string, cutoff time.Time) plotter.XYs {
	xys := plotter.XYs{}
	// xys := make(plotter.XYs, len(ms))
	for _, m := range ms {
		if !m.Created.Valid || m.Created.Time.Before(cutoff) {
			continue
		}
		xy := struct{ X, Y float64 }{X: float64(m.Created.Time.Unix())}
		switch field {
		case "battery":
			//xys[i].Y = m.Battery.Float64
			xy.Y = m.Battery.Float64
		case "humidity":
			//xys[i].Y = m.Humidity.Float64
			xy.Y = m.Humidity.Float64
		case "humiditytemp":
			//xys[i].Y = m.HumidityTemp.Float64
			xy.Y = m.HumidityTemp.Float64
		case "ihumidity":
			//xys[i].Y = m.Ihumidity.Float64
			xy.Y = m.Ihumidity.Float64
		case "ihumiditytemp":
			//xys[i].Y = m.IhumidityTemp.Float64
			xy.Y = m.IhumidityTemp.Float64
		case "pressure":
			//xys[i].Y = m.Pressure.Float64
			xy.Y = m.Pressure.Float64
		case "pressuretemp":
			//xys[i].Y = m.PressureTemp.Float64
			xy.Y = m.PressureTemp.Float64
		case "temperature":
			//xys[i].Y = m.Temperature.Float64
			xy.Y = m.Temperature.Float64
		case "temperatureext":
			// xys[i].Y = m.TemperatureExt.Float64
			xy.Y = m.TemperatureExt.Float64
		case "vref":
			//xys[i].Y = m.Vref.Float64
			xy.Y = m.Vref.Float64
		default:
			panic(fmt.Errorf("No such conversion field %q", field))
		}
		xys = append(xys, xy)
	}
	return xys
}
