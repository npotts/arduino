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

#include "barometer.h"


void InitMPL3115A2() { /*Setup Barometer*/
  barometer.begin(); // Get sensor online
  barometer.setModeBarometer(); // Measure pressure in Pascals from 20 to 110 kPa
  barometer.setOversampleRate(7); // Set Oversample to the recommended 128
  barometer.enableEventFlags();
}

/*ReadMPL3115A2 writes out DHT11 values to a serial port. On errors, it writes out nothing*/
void ReadMPL3115A2(JsonObject& obj) {
  float reading = barometer.readPressure() / 100.00;
  if (reading != -999.0) {
    obj["pressure"] = double_with_n_digits(reading, 10);
  }

  reading = barometer.readTemp();
  if (reading != -999.0) {
    obj["pressureTemp"] = double_with_n_digits(reading, 10);
  }
}
