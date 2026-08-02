#include <Arduino.h>
#include <WiFi.h>
#include "config.h"
#include "relay.h"
#include "leds.h"
#include "api.h"

Relay powerRelay(config::kPowerRelayPin);
Relay resetRelay(config::kResetRelayPin);
Led powerLed(config::kPowerLedPin);

void setup() {
  Serial.begin(115200);
  WiFi.mode(WIFI_STA);
  WiFi.begin(config::kSsid, config::kPassword);
  WiFi.setHostname(config::kHostname);

  powerRelay.begin();
  resetRelay.begin();
  powerLed.begin();

  setupApi();
}

void loop() {
  if (WiFi.status() != WL_CONNECTED) {
    delay(1000);
    return;
  }
}
