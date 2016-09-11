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

#include <math.h>

struct measurement readP() {
  struct measurement rtn = {NAN, NAN};
  rtn.a = barometer.readPressure() / 100.00;
  rtn.b = barometer.readTemp();
  return rtn;
}

struct measurement readRH() {
  struct measurement rtn = {NAN, NAN};
  rtn.a = rh.readHumidity();
  if (rtn.a == 998) {
    rh.begin();
    rtn.a = NAN;
    return rtn;
  }
  rtn.b = rh.readTemperature();
  return rtn;
}

struct measurement readPhoton() {
  struct measurement rtn = {NAN, NAN};
  rtn.a = analogRead(LIGHT); //read the light sensor
  rtn.b = analogRead(REFERENCE_3V3); //read the reference 3.3V signal to get the full-scale value
  rtn.b = 3.3 / rtn.b; //scale to get full-scale reading value
  rtn.a *= rtn.b; //expand light sensor to full scale again.
  return rtn;
}

struct measurement readBatt() {
  struct measurement rtn = {NAN, NAN};
  rtn.a = analogRead(BATT); //read battery voltage
  rtn.b = analogRead(REFERENCE_3V3); //read the reference 3.3V signal to get the full-scale value
  rtn.b = 3.3 / rtn.b; //scale to get full-scale reading value
  rtn.a *= rtn.b * 4.90; // Battery voltage is in voltage divider: (3.9k+1k)/1k
  return rtn;
}

struct measurement readTemp() {
  struct measurement rtn = {thermometer1.temperature(), thermometer2.temperature()};
  return rtn;
}


uint32_t index = 0;
void json() {
  Serial.print("{\"index\":"); Serial.print(index++);
  printoutJson("pressure", "ptemp", readP);
  printoutJson("rh", "rhtemp", readRH);
//  printoutJson("vphoton", "ref", readPhoton);
//  printoutJson("battery", "vref", readBatt);
  printoutJson("temperature_a", "temperature_b", readTemp);
  Serial.println("}");
}

void csv() {
  Serial.print(index++);
  printoutCsv(readTemp);
  printoutCsv(readP);
  printoutCsv(readRH);
  printoutCsv(readPhoton);
  printoutCsv(readBatt);
  Serial.println("");
}

