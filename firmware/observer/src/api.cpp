#include <Arduino.h>
#include <ESPAsyncWebServer.h>
#include "config.h"

AsyncWebServer server(config::kApiPort);

void setupApi() {
  server.on("/health", HTTP_GET, [](AsyncWebServerRequest* request) {
    request->send(200, "application/json", "{\"status\":\"ok\"}");
  });

  server.on("/status", HTTP_GET, [](AsyncWebServerRequest* request) {
    request->send(200, "application/json", "{\"status\":\"online\",\"firmware\":\"0.1.0\",\"uptime\":0}");
  });

  server.on("/environment", HTTP_GET, [](AsyncWebServerRequest* request) {
    request->send(200, "application/json", "{\"temperature\":0.0,\"humidity\":0.0}");
  });

  server.onNotFound([](AsyncWebServerRequest* request) {
    request->send(404, "application/json", "{\"error\":\"not found\"}");
  });

  server.begin();
}
