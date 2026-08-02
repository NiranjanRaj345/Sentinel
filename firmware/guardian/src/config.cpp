#include "config.h"

namespace config {
String wifiSsid() { return String(kSsid); }
String wifiPassword() { return String(kPassword); }
String hostname() { return String(kHostname); }
uint16_t apiPort() { return kApiPort; }
}
