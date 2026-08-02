#include <Arduino.h>
#include <ESPAsyncWebServer.h>
#include "config.h"

AsyncWebServer server(config::kApiPort);

void setupApi() {
  server.on("/health", HTTP_GET, [](AsyncWebServerRequest* request) {
    request->send(200, "application/json", "{\"status\":\"ok\"}");
  });

  server.onNotFound([](AsyncWebServerRequest* request) {
    request->send(404, "application/json", "{\"error\":\"not found\"}");
  });

  server.begin();
}
