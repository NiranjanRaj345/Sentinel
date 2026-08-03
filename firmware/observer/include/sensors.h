#ifndef OBSERVER_SENSORS_H
#define OBSERVER_SENSORS_H

#include <stdint.h>

class Sensors {
 public:
  Sensors();

  void begin();
  float readTemperature();
  float readHumidity();

 private:
  uint8_t tempPin_;
  uint8_t humPin_;
};

#endif
