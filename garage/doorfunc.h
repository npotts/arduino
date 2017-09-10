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

