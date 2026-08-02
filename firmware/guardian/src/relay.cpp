#include "config.h"

class Relay {
 public:
  Relay(uint8_t pin, bool activeLow = true)
      : pin_(pin), activeLow_(activeLow), state_(false) {}

  void begin() {
    pinMode(pin_, OUTPUT);
    write(false);
  }

  void write(bool on) {
    state_ = on;
    digitalWrite(pin_, activeLow_ ? (on ? LOW : HIGH) : (on ? HIGH : LOW));
  }

  void pulse(unsigned long ms = config::kRelayPulseMs) {
    write(true);
    delay(ms);
    write(false);
  }

  bool state() const { return state_; }

 private:
  uint8_t pin_;
  bool activeLow_;
  bool state_;
};
