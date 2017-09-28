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
#include "doorfunc.h"
#include "cmdprocessor.h"
/*
This is a really stupid sketch that monitors the A/D which is connected to a pot,
and a simple IO line that drives a relay. based on some commands.

*/

void setup() {
  Serial.begin(115200);
  init_fixup();
  analogReference(DEFAULT);
  pinMode(RelayPin, OUTPUT);
  digitalWrite(RelayPin, RelayIdle);
  updatePos();
  Serial.println(Prompt);
}

bool okfunc() { return false; }
void stopfunc() { Serial.print("-"); }
void upfunc() { Serial.print("/"); }
void dwfunc() { Serial.print("\\"); }
void testfunc() { optimize(0, 90000, &okfunc, &stopfunc, &upfunc, &dwfunc); }

void loop() {
  updatePos();
  cmd_proc();
}
