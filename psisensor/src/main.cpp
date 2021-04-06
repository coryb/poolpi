#include <Arduino.h>
#include <ArduinoJson.h>
#include <ArduinoLog.h>
#include <FS.h>
#include <HTTPClient.h>
#include <LITTLEFS.h>
#include <NTPClient.h>
#include <WiFi.h>
#include <driver/adc.h>
#include <esp_bt.h>
#include <esp_wifi.h>
#include <mem.h>

#define WIFI_TIMEOUT 30000  // 10 seconds
#define VOLTAGE_CALIBRATION 7.09

#define SLEEP10 (10L * 1000000L)
#define SLEEP1 (1L * 1000000L)

#define REPORT_DELAY (30)
#define TIMESYNC_DELAY (60 * 60 * 24)

#define SSID_MAX_LEN 32
#define WIFI_PW_MAX_LEN 63
#define INFLUXDB_MAX_LEN 512

RTC_DATA_ATTR uint64_t lastReportTime = 0;
RTC_DATA_ATTR uint64_t lastTimeSync = 0;
RTC_DATA_ATTR uint64_t lastTs = 0;
RTC_DATA_ATTR bool configLoaded = false;

WiFiUDP ntpUDP;
NTPClient timeClient(ntpUDP);
HTTPClient http;

RTC_DATA_ATTR char SSID[SSID_MAX_LEN];
RTC_DATA_ATTR char PASSWORD[WIFI_PW_MAX_LEN];
RTC_DATA_ATTR char INFLUXDB[INFLUXDB_MAX_LEN];

enum Mode { DEBUG_MODE,
            NORMAL_MODE };

Mode mode;

float readBattery();

struct timeval lastReadTV;

uint64_t nsTime() {
  if (lastReadTV.tv_sec == 0) {
    gettimeofday(&lastReadTV, NULL);
  }
  return lastReadTV.tv_sec * 1000000000LL + lastReadTV.tv_usec * 1000LL;
}

time_t sTime() {
  if (lastReadTV.tv_sec == 0) {
    gettimeofday(&lastReadTV, NULL);
  }
  return lastReadTV.tv_sec;
}

void loadConfig() {
  if (configLoaded) return;
  File file = LITTLEFS.open("/config.json");
  uint8_t buf[512];
  file.read(buf, 512);
  file.close();
  DynamicJsonDocument doc(512);
  deserializeJson(doc, buf);

  strncpy(SSID, doc.as<JsonObject>()["ssid"].as<const char*>(), SSID_MAX_LEN);
  strncpy(PASSWORD, doc.as<JsonObject>()["password"].as<const char*>(), WIFI_PW_MAX_LEN);
  strncpy(INFLUXDB, doc.as<JsonObject>()["influxdb"].as<const char*>(), INFLUXDB_MAX_LEN);
  configLoaded = true;
}

void blink() {
  if (mode == NORMAL_MODE) return;
  digitalWrite(LED_BUILTIN, HIGH);
  delay(100);
  digitalWrite(LED_BUILTIN, LOW);
  delay(100);
}

void deepSleep(uint64_t duration) {
  Log.verbose(F("Going to sleep..."));
  WiFi.disconnect(true);
  WiFi.mode(WIFI_OFF);
  btStop();

  adc_power_off();
  esp_wifi_stop();
  esp_bt_controller_disable();

  blink();
  // Configure the timer to wake us up!
  esp_sleep_enable_timer_wakeup(duration);

  // Go to sleep! Zzzz
  esp_deep_sleep_start();
}

void wifiConnect() {
  WiFi.begin(SSID, PASSWORD);
  Log.verbose(F("Connecting to WiFi..."));
  unsigned long s = millis();

  // Keep looping while we're not connected AND haven't reached the timeout
  while (WiFi.status() != WL_CONNECTED && millis() - s < WIFI_TIMEOUT) {
    delay(10);
  }

  // Make sure that we're actually connected, otherwise go to deep sleep
  if (WiFi.status() != WL_CONNECTED) {
    Log.error(F("FAILED"));
    deepSleep(SLEEP10);
  }
  Log.verbose(F("Connected"));
}

void syncTime() {
  timeClient.forceUpdate();
  timeval tv;
  tv.tv_sec = timeClient.getEpochTime();
  settimeofday(&tv, NULL);
  lastTimeSync = sTime();
}

char buf[1024];

void record(fs::File& file, String name, float value) {
  int len = sprintf(buf, "%s value=%.2f %llu\n", name.c_str(), value, nsTime());
  file.write((uint8_t*)buf, len);
  Log.verbose(String(buf).substring(0, len - 1).c_str());
}

void printTimestamp(Print* _logOutput) {
  time_t ts = sTime();
  tm* t = localtime(&ts);
  char c[21];
  strftime(c, 21, "%Y.%m.%d-%H:%M:%S ", t);
  _logOutput->print(c);
}

void printNewline(Print* _logOutput) {
  _logOutput->print('\n');
}

void setup() {
  setCpuFrequencyMhz(80);
  Serial.begin(9600);
  Serial.println();
  pinMode(LED_BUILTIN, OUTPUT);
  pinMode(GPIO_NUM_16, INPUT);
  pinMode(GPIO_NUM_27, INPUT_PULLUP);

  mode = (Mode)digitalRead(GPIO_NUM_27);
  int logLevel = LOG_LEVEL_ERROR;
  if (mode == DEBUG_MODE) {
    Log.setPrefix(printTimestamp);
    Log.setSuffix(printNewline);
    logLevel = LOG_LEVEL_VERBOSE;
  }
  Log.begin(logLevel, &Serial);
  Log.verbose(F("Starting"));

  if (!LITTLEFS.begin(false)) {
    Log.error(F("LITTLEFS Mount Failed"));
    return;
  }

  uint64_t ts = sTime();
  if (ts < 5000) {
    Log.notice(F("Time not set, doing so now"));
    loadConfig();
    wifiConnect();
    syncTime();
    // we are connected to WiFi now which will interfere with out
    // ADC reads, so just sleep this first time.
    deepSleep(SLEEP1);
  }

  if (lastReportTime == 0) {
    Log.notice(F("Resetting lastReportTime"));
    lastReportTime = ts;
  }

  File report = LITTLEFS.open("/report.txt", FILE_APPEND);

  adc_power_on();
  uint16_t psi = 0;
  for (uint8_t i = 0; i < 10; i++) {
    psi += analogRead(GPIO_NUM_34);
  }

  record(report, F("psi"), float(psi / 10));

  // first reading might be bogus if we used wifi recently, so read twice
  float battery = readBattery();
  battery = readBattery();
  record(report, F("battery"), battery);

  if (lastReportTime + REPORT_DELAY > ts) {
    report.close();
    deepSleep(SLEEP10);
  }

  loadConfig();
  wifiConnect();

  record(report, F("wifiLevel"), WiFi.RSSI());

  report = LITTLEFS.open("/report.txt", FILE_READ);
  size_t size = report.size();
  uint8_t* reportBuf = new uint8_t[size];
  report.read(reportBuf, size);
  report.close();

  http.begin(INFLUXDB);
  int resp = http.POST(reportBuf, size);
  http.end();
  Log.verbose(F("POST %d"), resp);
  delete[] reportBuf;

  if ((resp >= 200 && resp < 300) || resp == 400) {
    // truncate report
    File report = LITTLEFS.open("/report.txt", FILE_WRITE);
    report.close();
    lastReportTime = ts;
    if (resp == 400) {
      blink();
      blink();
    } else {
      blink();
    }
  } else {
    blink();
    blink();
    blink();
  }

  if (lastTimeSync + TIMESYNC_DELAY > ts) {
    syncTime();
  }

  deepSleep(SLEEP10);
}

void loop() {
}

// readBattery from:
// https://github.com/G6EJD/LiPo_Battery_Capacity_Estimator/blob/master/ReadBatteryCapacity_LIPO.ino
/* An improved battery estimation function 
   This software, the ideas and concepts is Copyright (c) David Bird 2019 and beyond.
   All rights to this software are reserved.
   It is prohibited to redistribute or reproduce of any part or all of the software contents in any form other than the following:
   1. You may print or download to a local hard disk extracts for your personal and non-commercial use only.
   2. You may copy the content to individual third parties for their personal use, but only if you acknowledge
      the author David Bird as the source of the material.
   3. You may not, except with my express written permission, distribute or commercially exploit the content.
   4. You may not transmit it or store it in any other website or other form of electronic retrieval system for commercial purposes.
   5. You MUST include all of this copyright and permission notice ('as annotated') and this shall be included in all copies
      or substantial portions of the software and where the software use is visible to an end-user.
   THE SOFTWARE IS PROVIDED "AS IS" FOR PRIVATE USE ONLY, IT IS NOT FOR COMMERCIAL USE IN WHOLE OR PART OR CONCEPT.
   FOR PERSONAL USE IT IS SUPPLIED WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
   WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
   IN NO EVENT SHALL THE AUTHOR OR COPYRIGHT HOLDER BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN
   AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
   OTHER DEALINGS IN THE SOFTWARE.
*/
float readBattery() {
  float percentage = 100;
  float voltage = analogRead(GPIO_NUM_4) / 4096.0 * VOLTAGE_CALIBRATION;
  Log.verbose(F("Voltage %F"), voltage);
  percentage = 2808.3808 * pow(voltage, 4) - 43560.9157 * pow(voltage, 3) +
               252848.5888 * pow(voltage, 2) - 650767.4615 * voltage +
               626532.5703;
  if (voltage > 4.19)
    percentage = 100;
  else if (voltage <= 3.50)
    percentage = 0;
  return percentage;
}
