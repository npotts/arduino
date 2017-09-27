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

