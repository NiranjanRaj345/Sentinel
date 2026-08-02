#include "config.h"

class Led {
 public:
  explicit Led(uint8_t pin) : pin_(pin), state_(false) {}

  void begin() {
    pinMode(pin_, INPUT);
    state_ = false;
  }

  bool read() const {
    return digitalRead(pin_) == HIGH;
  }

 private:
  uint8_t pin_;
  bool state_;
};
