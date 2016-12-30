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

#include "1wire-temp.h"


/*Search for 1wire devices on 1w bus*/
bool W1temp::locate(unsigned int skip) {
  int mskip = skip;
  thermo.reset_search();
  thermo.target_search(family); //temp family
  bool ok;
  
  while (true) { //prod through
    if (!thermo.search(w1addr)) {
      // Serial.println("-- Unable to search for devices");
      return false;
    }
    if (skip == 0) {
      break;
    }
    // Serial.print("-- Skipping #"); Serial.print(skip); Serial.print(" ");
    for(byte i = 0; i<8; i++) {
      // Serial.print(w1addr[i], HEX);
    }
    // Serial.println();
    skip--;
  }

//  if (!thermo.search(w1addr)) { return false; }
  // Serial.print("-- Found sensor at ");
  for(byte i = 0; i<8; i++) {
    // Serial.print(w1addr[i], HEX);
  }
  if (OneWire::crc8(w1addr, 7) != w1addr[7]) {
      // Serial.println(". CRC invalid!");
      return false;
  }
  // Serial.println("");
  return true;
}
 
bool W1temp::measure(void) {
  if (!thermo.reset()) {return false;};
  thermo.select(w1addr);
  thermo.write(0x44); //self-powered
  tstart = ::millis();
}

//temperature returns the temperature, ready is set to true if the measurement is valid
float W1temp::temperature() {
  //  while(::millis() - tstart < millis); //wait at least millis - DETECT CLOCK ROLLOVER
  if (!thermo.reset()) {return NAN;} //wake up device - should always get a pulse
  thermo.select(w1addr);
  thermo.write(0xBE); //read scratchpad
  
  byte data[12];
  for (byte i = 0; i < 9; i++) { data[i] = thermo.read(); } // we need 9 bytes 
  if (OneWire::crc8(data, 8) != data[8]) { return false; }
  return tempC(data); 
}

float W1temp::tempC(byte data[12]) {
  int16_t raw = (data[1] << 8) | data[0];
  byte cfg = (data[4] & 0x60);
  // at lower res, the low bits are undefined, so let's zero them
  if (cfg == 0x00) raw = raw & ~7;  // 9 bit resolution, 93.75 ms
  else if (cfg == 0x20) raw = raw & ~3; // 10 bit res, 187.5 ms
  else if (cfg == 0x40) raw = raw & ~1; // 11 bit res, 375 ms
  //// default is 12 bit resolution, 750 ms conversion time
  return (float) raw / 16.0;
}

