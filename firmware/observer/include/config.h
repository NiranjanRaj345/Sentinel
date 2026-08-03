#ifndef OBSERVER_CONFIG_H
#define OBSERVER_CONFIG_H

#include <Arduino.h>

namespace config {
const char* kSsid = "";
const char* kPassword = "";
const char* kHostname = "observer";
const uint16_t kApiPort = 80;

const uint8_t kSdaPin = 21;
const uint8_t kSclPin = 22;

const unsigned long kStatusPollMs = 1000;
}

#endif
