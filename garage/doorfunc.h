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
#include "hwcfg.h"

//Declaration of types of door motion
enum Motion {
  Stopped, //Not moving
  MovingUp, //Door is raising
  MovingDown //Door is falling
};


unsigned int pos; //wiper position

//forward declarations
void triggerRelay();
Motion DoorMotion();
void updatePos();
bool same(unsigned int, unsigned int);
bool isClosed();
bool wideOpen();
bool isOpen();
bool optimize(int, unsigned long, bool (*)(void), void (*)(void), void (*)(void), void (*)(void) );
bool openDoor();
bool closeDoor();
int percentOpen(unsigned int);

unsigned int averagePosition(unsigned char );
