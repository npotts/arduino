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

#include <Wire.h>
#include <ArduinoJson.h>

#include "DHT11.h"
#include "barometer.h"
#include "htu21d.h"
#include "1wire-temp.h"

const byte REFERENCE_3V3 = A3;
const byte BATT = A2;


/*both temps are on same bus, same family*/
W1temp thermometer1(10, 0x28, 0, 800);
W1temp thermometer2(10, 0x28, 0, 800);

/*tmeasure starts measurement of the values*/
void tmeasure() {
  thermometer1.measure();
  thermometer2.measure();
}

/*ReadTemps fetches the values from the temp sensors*/
void ReadTemps(JsonObject& obj) {
  obj["temperatureExt"] = double_with_n_digits(thermometer1.temperature(),10);
  obj["temperature"] = double_with_n_digits(thermometer2.temperature(), 10);
} 

void ReadSystem(JsonObject& obj) {
  float vref = 3.3 / analogRead(REFERENCE_3V3); //read the reference 3.3V signal to get the full-scale value
  float battery = analogRead(BATT) * vref * 4.90; // Battery voltage is in voltage divider: (3.9k+1k)/1k
  obj["vref"] = double_with_n_digits(vref, 5);
  obj["battery"] = double_with_n_digits(battery, 5);
}

void setup() {
  pinMode(REFERENCE_3V3, INPUT);
  pinMode(BATT, INPUT);
  Serial.begin(115200);
  InitMPL3115A2();
  InitHTU21D();
  while(!thermometer1.locate(0)) {}
  while(!thermometer2.locate(1)) {}
}

void waitFor(long int start, long int atleast) {
  while (1) { //wait full 800ms:
    long int now = millis();
    if (now < start) { return; } // clock rollover.  Reboot
    if (now - start > atleast) return;
  }
}

void loop() {
  StaticJsonBuffer<256> jsonBuffer; //json buffer
  JsonObject& object = jsonBuffer.createObject();
  long int start = millis();

  tmeasure(); //start temp measuring
  ReadSystem(object);
  ReadDH11(object);
  ReadMPL3115A2(object);
  ReadHTU21D(object);
  waitFor(start, 800); //wait at least 800ms
  ReadTemps(object);

  object.printTo(Serial);
  Serial.println();
  waitFor(start, 5000); //wait at least 5000ms
}

