/*
The MIT License (MIT)

Copyright (c) 2016-2017 Nick Potts

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

#include "hwcfg.h"

/*use EEPROM, addres EEPROMSTART {+0,+1} is Floor, EEPROMSTART {+2,+3} is the ceiling*/
void init_readPositionValues() {
  DoorAtFloor = (EEPROM.read(EEPROMSTART+0) << 8) + EEPROM.read(EEPROMSTART+1);;
  DoorAtCeiling = (EEPROM.read(EEPROMSTART+2) << 8) + EEPROM.read(EEPROMSTART+3);
}

/*use EEPROM, addres EEPROMSTART {+0,+1} is Floor, EEPROMSTART {+2,+3} is the ceiling*/
void writePositionValues(unsigned int floor, unsigned int ceiling) {
  EEPROM.write(EEPROMSTART, (unsigned char)(floor >> 8));
  EEPROM.write(EEPROMSTART + 1, (unsigned char)(floor & 0xFF));
  EEPROM.write(EEPROMSTART + 2, (unsigned char)(ceiling >> 8));
  EEPROM.write(EEPROMSTART + 3, (unsigned char)(ceiling & 0xFF));  
}


//fixup should alter pdelta and PosIncreaseOpensDoor
unsigned int init_fixup() {
  init_readPositionValues();
  PosIncreaseOpensDoor = DoorAtCeiling > DoorAtFloor;
}
