#ifndef GUARDIAN_CONFIG_H
#define GUARDIAN_CONFIG_H

#include <Arduino.h>

namespace config {
const char* kSsid = "";
const char* kPassword = "";
const char* kHostname = "guardian";
const uint16_t kApiPort = 80;

const uint8_t kPowerRelayPin = 26;
const uint8_t kResetRelayPin = 27;
const uint8_t kPowerLedPin = 34;

const unsigned long kRelayPulseMs = 200;
const unsigned long kStatusPollMs = 1000;
}

#endif
