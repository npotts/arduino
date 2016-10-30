/*
 Fading

 This example shows how to fade an LED using the analogWrite() function.

 The circuit:
 * LED attached from digital pin 9 to ground.

 Created 1 Nov 2008
 By David A. Mellis
 modified 30 Aug 2011
 By Tom Igoe

 http://www.arduino.cc/en/Tutorial/Fading

 This example code is in the public domain.

 */

 #include "MotorCtrl.h"

 MotorCtrl left(6,7);
 MotorCtrl right(5,4);

void setup() {
  Serial.begin(115200);
}

void loop() {
  left.forward();
  right.forward();
  delay(1000);

  left.speed(255);
  right.speed(255);
  delay(1000);

  left.back();
  right.back();
  delay(1000);

  left.speed(127);
  right.speed(127);
  delay(1000);

}