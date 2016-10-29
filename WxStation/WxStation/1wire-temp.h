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


#pragma once
#define ONEWIRE_CRC8_TABLE 1 //Use fast CRC Lookup table
#include "OneWire.h"

class W1temp {
private:
    int millis, twowire;
    uint8_t pin, family;
    byte w1addr[8]; //temp address
    OneWire thermo;
    uint16_t tstart; //this will eventually roll over...
    float tempC(byte data[12]);
    
public:
    /*Initiate a 1Wire temperature sensor attached to pin, as a device in the specified family. If twowire  
     * is 0 then the temp sensor should be powered on permamently via external vcc to 3.3-5.5VDC.  min_millis
     * is the length of time to wait until reading values off the device.  Prone to rollover (for now)*/
    W1temp(uint8_t pin, uint8_t family, int twowire=0, int min_millis=750): pin(pin), family(family), millis(min_millis), thermo(pin), twowire(twowire) {}
    
    //locate searches for a W1temp item, and returns true if one is found
    bool locate(unsigned int skip=0); //returns true if a sensor from family is found

    //measure starts measuring a temperature - returning false if it encountered errors (eg, 1wire device removed)
    bool measure(void);

    /*temperature returns true if the measurement worked, and the value written to temperature. 
     * if false is written, you need to call measure before you recall temperature() */
    float temperature();
};

