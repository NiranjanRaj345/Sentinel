#ifndef OBSERVER_DISPLAY_H
#define OBSERVER_DISPLAY_H

#include <stdint.h>
#include "config.h"

class Display {
 public:
  Display(uint8_t sda = config::kSdaPin, uint8_t scl = config::kSclPin);

  void begin();
  void showStatus(const char* status, const char* firmware);

 private:
  uint8_t sda_;
  uint8_t scl_;
};

#endif
