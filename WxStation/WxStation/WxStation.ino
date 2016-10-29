#include <ArduinoSTL.h>
#include <Wire.h>
#include <ArduinoJson.h>

#include "DHT11.h"
#include "barometer.h"
#include "htu21d.h"

const byte REFERENCE_3V3 = A3;
const byte BATT = A2;

StaticJsonBuffer<256> jsonBuffer; //json buffer

void system(JsonObject& obj) {
  float vref = 3.3 / analogRead(REFERENCE_3V3); //read the reference 3.3V signal to get the full-scale value
  float battery = analogRead(BATT) * vref * 4.90; // Battery voltage is in voltage divider: (3.9k+1k)/1k
  obj["vref"] = vref;
  obj["battery"] = battery;
}

void setup() {
  pinMode(REFERENCE_3V3, INPUT);
  pinMode(BATT, INPUT);
  Serial.begin(115200);
  InitMPL3115A2();
  InitHTU21D();
}

void loop() {
  JsonObject& object = jsonBuffer.createObject();
  ReadDH11(object);
  ReadMPL3115A2(object);
  ReadHTU21D(object);
  object.printTo(Serial);
  Serial.println();
  delay(5000);
}

