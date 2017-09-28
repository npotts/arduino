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

#include "cmdprocessor.h"

/*cmd_status returns the full status which is a comma
seperated string terminated with a newline of the
following parameters:
 - millis() when called
 - pos: raw A/D count
 - DoorAtFloor: A/D count when at the floow
 - DoorAtCeiling: A/D count when at the ceiling
 - Percentage (0 to 100) open.  Open = 100%
 - Door Is at floor (true | false)
 - Door Is at ceiling (true | false)
*/
void cmd_status() {
  Serial.print(millis()); Serial.print(",");
  Serial.print(pos); Serial.print(",");
  Serial.print(DoorAtFloor); Serial.print(",");
  Serial.print(DoorAtCeiling); Serial.print(",");
  Serial.print(percentOpen(pos)); Serial.print(",");
  Serial.print(isClosed()); Serial.print(",");
  Serial.println(wideOpen());
}


void cmd_help() {
  Serial.print(F("This is the idiotic garage door opener. I speak the following ASCII\n"
  "?\trequest the following CSV formed info:\n"
  "\t\tmillis\n"
  "\t\tpos raw A/D count\n"
  "\t\tA/D count when door is at floor\n"
  "\t\tA/D count when door is at the ceiling\n"
  "\t\t%open.  Open=100%\n"
  "\t\tDoor at floor: 0|1\n"
  "\t\tDoor at ceiling: 0|1\n"
  "^\tTrigger the relay\n"
  "h\tThis limited 'help'\n"
  "c\tAttempt to Close Door\n"
  "o\tAttempt to Open Door\n"
  "m\tDoor Movement info\n"
  "~\tPerform a Door-Calibration.  Make sure the door is closed and wait.  This could take a long time.  | symbols are used to detect motion, * means stopped.  AFter cal is complete, it writes values to EEPROM and resumes normal operations\n"  
  ));
}

/*cmd_proc takes one of the following single ASCII character
incoming 'commands' and does something with them. It explicitly
clears the queue after enacting on one.

ASCII  Response
?      request position information
^      trigger the relay
h      Some limited help
c      Close Door - attempts to close the door.  This uses A2D values to button smash until the door is shut
o      Open Door -- same as above, but in reverse.
m      Door Moving data: 0=stationary, 1=Door Moving Up, 2=Door Moving Down
*/
void cmd_proc() {
  bool cont = true;
  while(cont && Serial.available() > 0) {
    switch(Serial.read()) {
      case '?':
        cont = false;
        cmd_status();
        break;
      case '^':
        cont = false;
        Serial.println("Pushing button");
        triggerRelay();
        break;
      case 'c':
        cont = false;
        Serial.print("Attemping to close door: ");
        Serial.println(closeDoor() ? "closed": "nope");
        break;
      case 'o':
        cont = false;
        Serial.print("Attemping to open door: ");
        Serial.println(openDoor() ? "opened": "nope");
        break;
      case 'm':
        cont = false;
        Serial.print("Door is currently ");
        switch(DoorMotion()) {
          case Stopped: Serial.println("stopped"); break;
          case MovingUp: Serial.println("moving up"); break;
          case MovingDown: Serial.println("moving down"); break;
        }
        break;
      case '~':
        cmd_calibrate();
        break;
      case 'h':
        cont = false;
        cmd_help();
        break;
      default:
        continue;
        break;
    }
    Serial.println(Prompt);
  }
}

/*cmd_calibrate performs a calibration of the garage door.  It expects 
the door to be closed initially, and then does the following:

- Measures 16 samoples at the currrent position (closed)
- Smashes the button
- Waits until the M command responds with "stopped"  5 times
- Measures 16 samples at the currrent position (open)
- Writes values to the EEPROM
- Resets Device


*/
void cmd_calibrate() {
  int count = 0;
  Serial.print("Cal> ");
  unsigned int floor = averagePosition(4);
  Serial.print("Floor:"); Serial.print(floor);
  triggerRelay(); delay(2000); //trigger button and wait 2s for door to start
  Serial.print(" {");
  while (true) {
    switch(DoorMotion()) {
      case Stopped:
        Serial.print("*");
        count++;
        break;
      default:
        Serial.print("|");
        count = 0;
        break;
    }
    delay(300); //wait another 200ms
    if (count >= 5) break; //motor isnt moving
  }
  Serial.print("} Stopped. ");
  unsigned int ceiling = averagePosition(4);
  Serial.print("Ceiling:"); Serial.print(ceiling);
  writePositionValues(floor, ceiling);
  init_readPositionValues();
  Serial.print(" Done\n");
}

