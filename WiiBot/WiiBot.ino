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
