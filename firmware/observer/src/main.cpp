#include <Arduino.h>
#include <WiFi.h>
#include "config.h"
#include "display.h"
#include "sensors.h"
#include "api.h"

Display display;
Sensors sensors;

void setup() {
  Serial.begin(115200);
  WiFi.mode(WIFI_STA);
  WiFi.begin(config::kSsid, config::kPassword);
  WiFi.setHostname(config::kHostname);

  display.begin();
  sensors.begin();

  setupApi();
}

void loop() {
  if (WiFi.status() != WL_CONNECTED) {
    delay(1000);
    return;
  }
}
