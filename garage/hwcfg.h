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

#pragma once
/*
# Diagram

DoorAtFloor and DoorAtCeiling instruct this firmware the range the door should operate in, as well as which way is 'up'

--Ceiling---|  <--- DoorAtCeiling value
            |
  Garage    |   Door Operating range
  or inside |
            |
---Ground---|  <--- DoorAtFloow value

*/
#define DoorAtFloor 165 //A2D values lower than this mean the door is closed
#define DoorAtCeiling 878 //A2D values greater than this mean the door is closed
#define A2DJitter 5 // How many counts +/- from a value do we consider as 'the same'.  THis is effectively your noise

//Don't alter the following - It is used to pick the right strategy
#if DoorAtCeiling > DoorAtFloor
#define PosIncreaseOpensDoor true
unsigned int pdelta = DoorAtCeiling - DoorAtFloor;
#else
#define PosIncreaseOpensDoor false
unsigned int pdelta = DoorAtFloor - DoorAtCeiling;
#endif

//Pin Values
#define A2DPin A0
#define RelayPin 13
#define RelayDur 300  //300ms
#define RelayIdle LOW //Idling should be low
#define RelaySet HIGH //when set, should set the pin high

//Motion algo tuneables
#define MotionWait 300      // How long to wait between position samples

#define Prompt ">>"
