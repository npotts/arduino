#include "doorfunc.h"

/*triggerRelay triggers the relay for some time*/
void triggerRelay() {
  digitalWrite(RelayPin, RelaySet);
  delay(RelayDur);
  digitalWrite(RelayPin, RelayIdle);
}

/*updatePos updates the global pos variable*/
void updatePos() {
  pos = analogRead(A2DPin);
} 

/*same returns true if a and b are both within A2Djitter */
bool same(unsigned int a, unsigned int b) {
  return (a > b) ? (a - b < A2DJitter) : (b - a < A2DJitter);
}

/*isClosed returns true if the door is in the 'closed' position*/
bool isClosed() {
  return same(pos, DoorAtFloor) || 
    (PosIncreaseOpensDoor && pos < DoorAtFloor) ||
    (!PosIncreaseOpensDoor && pos > DoorAtFloor);
}

/*wideOpen returns true if the door is in the fully opened position*/
bool wideOpen() {
  return same(pos, DoorAtCeiling) ||
    (PosIncreaseOpensDoor && pos > DoorAtCeiling) ||
    (!PosIncreaseOpensDoor && pos < DoorAtCeiling);
}

/*isOpen returns true if the door is in any sort of 'open' position*/
bool isOpen() {
  return !isClosed();
}

/*DoorMotion returns the guessed Motion of the door: Up, Down, or stopped.
This takes ~300ms to perform, and assumes that if the A2D counts are not moving by more than 10 counts,
it isnt moving*/
Motion DoorMotion() {
  unsigned int then = pos; //capture value now
  delay(MotionWait);
  updatePos();
  bool increasing = (pos < then);
  if (same(then, pos)) return Stopped;
  if ( (PosIncreaseOpensDoor && increasing) || (!PosIncreaseOpensDoor && !increasing)) return MovingUp;
  return MovingDown;
}
/*This handy function takes some time periods and function pointers to attempt to put the door in some known position.
It loops up to duration numberof milliseconds, with a pause of at least loopdelay
*/
bool optimize(int loopdelay, unsigned long duration, bool (*exitCond)(void), void (*onstopped)(void), void (*onup)(void), void (*ondown)(void) ){
  unsigned long start = millis();
  while((unsigned long)(millis() - start) < duration) {
    updatePos();
    if (exitCond()) return true;
    switch(DoorMotion()) {
      case Stopped: // not moving
        onstopped();
        break;
       case MovingUp: //moving up
        onup();
        break;
       case MovingDown: //moving down
        ondown();
        break;
      default:
        return false;
    }
    delay(loopdelay); //wait for stuff to happen
  }
  return false;
}
bool openDoor() {
  //every 500ms, for 30000ms, exit if wideOpen, smashing the button if going down, or stopped
  return optimize(2000, 90000, &wideOpen, &triggerRelay, &updatePos, &triggerRelay);
}

bool closeDoor() {
  //every 500ms, for 30000ms, exit if isClosed(), smashing the button if going up, or stopped
  return optimize(2000, 90000, &isClosed, &triggerRelay, &triggerRelay, &updatePos);
}

/*this returns a number, from 0 to 100 that reperesents the open state of the door.
0 should be understood to be 100% closed
100 should be understood to be 100% open.

npos is the current ADC position of the door
*/
int percentOpen(unsigned int npos) {
  unsigned int partial;
  if (PosIncreaseOpensDoor) {
    if (npos < DoorAtFloor) return 0;
    if (npos > DoorAtCeiling) return 100;
    partial = 100 * (npos - DoorAtFloor);
  } else {
    if (npos > DoorAtFloor) return 0;
    if (npos < DoorAtCeiling) return 100;
    partial = 100 * (DoorAtFloor - npos);
  }
  partial /= pdelta;
  return partial;
}

