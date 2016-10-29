# Purpose
I want to know what the temperature is outside

# Schematic


# Parameters
The following parameters are recorded and stored.  Missing, or invalid, values are thrown away:

  - Pressure (barometric, mbar) & pressure sensors' temperature
  - Internal humidity & temperature
  - External Humidity & temperature
  - External Temperature
  - Roof Temperature

# Arduino output

The system makes a measurement roughly every 5 seconds, and the host OS pushes it into [brianiac](http://github.com/npotts/brianiac) for posterities sake.

## Format

A typical data frame emitted over the serial line looks something like this (with extra spaces, newlines, etc) (:
```js
  {
    "pressure_mbar": 843.5235,
    "pressure_temp": 25.3252,
    "ihumidity": 56.343,
    "ihumidity_temp": 25.232,
    "humidity": 60,
    "humidity_temp": 25,
    "temperature1": 35.2213,
    "temperature2": 35.2213,
  }
```

