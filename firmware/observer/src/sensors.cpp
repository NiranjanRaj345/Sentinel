#include "sensors.h"

Sensors::Sensors() : tempPin_(A0), humPin_(A1) {}

void Sensors::begin() {}

float Sensors::readTemperature() {
  return 0.0f;
}

float Sensors::readHumidity() {
  return 0.0f;
}
