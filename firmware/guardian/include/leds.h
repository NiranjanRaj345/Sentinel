#ifndef GUARDIAN_LEDS_H
#define GUARDIAN_LEDS_H

#include <stdint.h>

class Led {
 public:
  explicit Led(uint8_t pin);
  void begin();
  bool read() const;

 private:
  uint8_t pin_;
  bool state_;
};

#endif
