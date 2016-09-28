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

#include "measurements.h"


struct frame fetch() {
  struct frame rtn = {index++, NAN, NAN, NAN, NAN, NAN, NAN, NAN, NAN};
  thermometer1.measure(); //prep read
  thermometer2.measure(); //prep read
  rtn.pressure = barometer.readPressure() / 100.00; if (rtn.pressure < 0) { rtn.pressure = NAN; }
  rtn.ptemp = barometer.readTemp(); if (rtn.ptemp == -999.0) { rtn.ptemp = NAN; }
  rtn.humidity = rh.readHumidity();
  if (rtn.humidity == 998) { rh.begin(); rtn.humidity = NAN; } else { rtn.htemp = rh.readTemperature(); }
  rtn.vref = 3.3 / analogRead(REFERENCE_3V3); //read the reference 3.3V signal to get the full-scale value
  rtn.battery = analogRead(BATT) * rtn.vref * 4.90; // Battery voltage is in voltage divider: (3.9k+1k)/1k
  //wait full 800ms:
  delay(600);
  
  rtn.tempA = thermometer1.temperature();
  rtn.tempB = thermometer2.temperature();
  return rtn;
}

void niceJson(String key, float val) {
  if (val == val) {
    Serial.print(",\""); Serial.print(key); Serial.print("\":"); Serial.print(val);
  }
}

void niceCsv(float val) {
  Serial.print(",");
  if (val == val) {
    Serial.print(val);
  }
}

void niceSql(float val) {
  Serial.print(",");
  if (val == val) {
    Serial.print(val, 6);
  } else {
    Serial.print("NULL");
  }
}

void sqlinsert(struct frame frame) {
  Serial.print("INSERT INTO %s (indx, pressure, tempa, tempb, humidity, ptemp, htemp, battery) VALUES (");
  Serial.print(frame.index);
  niceSql(frame.pressure);
  niceSql(frame.tempA);
  niceSql(frame.tempB);
  niceSql(frame.humidity);
  niceSql(frame.ptemp);
  niceSql(frame.htemp);
  niceSql(frame.battery);
  Serial.println(");");
}

void json(struct frame frame) {
  Serial.print("{\"index\":"); Serial.print(frame.index);
  niceJson("pressure", frame.pressure);
  niceJson("temp_a", frame.tempA);
  niceJson("temp_b", frame.tempB);
  niceJson("humidity", frame.humidity);
  niceJson("ptemp", frame.ptemp);
  niceJson("htemp", frame.htemp);
  niceJson("battery", frame.battery);
  Serial.println("}");
}


