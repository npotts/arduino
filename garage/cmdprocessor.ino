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
  Serial.print("This is the idiotic garage door opener. I speak the following ASCII\n"
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
  );
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

