#include "hwcfg.h"
#include "doorfunc.h"
#include "cmdprocessor.h"
/*
This is a really stupid sketch that monitors the A/D which is connected to a pot, 
and a simple IO line that drives a relay. based on some commands.

*/

void setup() {
  Serial.begin(115200);
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
