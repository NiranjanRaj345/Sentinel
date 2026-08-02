#ifndef GUARDIAN_RELAY_H
#define GUARDIAN_RELAY_H

#include <stdint.h>
#include "config.h"

class Relay {
 public:
  Relay(uint8_t pin, bool activeLow = true);
  void begin();
  void write(bool on);
  void pulse(unsigned long ms = config::kRelayPulseMs);
  bool state() const;

 private:
  uint8_t pin_;
  bool activeLow_;
  bool state_;
};

#endif
