#include "display.h"

Display::Display(uint8_t sda, uint8_t scl) : sda_(sda), scl_(scl) {}

void Display::begin() {
  // Placeholder for OLED/Screen initialization
}

void Display::showStatus(const char* status, const char* firmware) {
  (void)status;
  (void)firmware;
}
