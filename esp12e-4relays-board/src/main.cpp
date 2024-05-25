#include <Arduino.h>
#include <ESP8266WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>
#include "credentials.h"
#include <string.h>

#define R1 15
#define R2 14
#define R3 12
#define R4 13

WiFiClient espClient;
PubSubClient client(espClient);
unsigned long lastMsg = 0;
#define MSG_BUFFER_SIZE  (50)
char msg[MSG_BUFFER_SIZE];
int value = 0;

String random_string = String(random(0xffff), HEX);
String clientId = "ESP8266Client-";

void setup_wifi() {

  delay(10);
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  randomSeed(micros());

  Serial.println("");
  Serial.println("WiFi connected");
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());
}

void turn_on(int line_number){

  if(line_number == 1)
    digitalWrite(R1, HIGH);
  else if (line_number == 2)
    digitalWrite(R2, HIGH);
  else if (line_number == 3)
    digitalWrite(R3, HIGH);
  else if (line_number == 4)
    digitalWrite(R4, HIGH);
}

void turn_off(int line_number){

  if(line_number == 1)
    digitalWrite(R1, LOW);
  else if (line_number == 2)
        digitalWrite(R2, LOW);
  else if (line_number == 3)
    digitalWrite(R3, LOW);
  else if (line_number == 4)
    digitalWrite(R4, LOW);
}

void callback(char* topic, byte* payload, unsigned int length) {
  Serial.print("Message arrived [");
  Serial.print(topic);
  Serial.print("] ");
  for (int i = 0; i < length; i++) {
    Serial.print((char)payload[i]);
  }
  Serial.println();

  StaticJsonDocument<200> doc;
  DeserializationError error = deserializeJson(doc, payload, length);

  if (error) {
    Serial.print(F("deserializeJson() failed: "));
    Serial.println(error.f_str());
    return;
  }

  /*
  **  We're trying to parse this kind of payload:
  **  {"line": 1, "cmd": "ON"}
  */

  int line_number = doc["line"];
  const char* cmd = doc["cmd"];

  Serial.println(line_number);
  Serial.println(cmd);

  if (strcmp(cmd, "ON") == 0) {
    turn_on(line_number);
  } else if (strcmp(cmd, "OFF") == 0) {
    turn_off(line_number);
  } else {
    Serial.println("Command not recognized: " + String(cmd));
  }

}

void reconnect() {
  // Loop until we're reconnected
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    // Create a random client ID
    clientId += random_string;
    // Attempt to connect
    if (client.connect(clientId.c_str())) {
      Serial.println("connected");
      // Once connected, publish an announcement...
      client.publish(service_topic, clientId.c_str());
      // ... and resubscribe
      client.subscribe(lines_topic);
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      // Wait 5 seconds before retrying
      delay(5000);
    }
  }
}


void setup() {
  pinMode(BUILTIN_LED, OUTPUT);
  Serial.begin(115200);
  setup_wifi();
  client.setServer(mqtt_server, mqtt_port);
  client.setCallback(callback);
  pinMode(R1, OUTPUT);
  pinMode(R2, OUTPUT);
  pinMode(R3, OUTPUT);
  pinMode(R4, OUTPUT);
  digitalWrite(2, HIGH); 

}

void blink(int times){
  for (int i=0; i<times; i++){
      digitalWrite(2, LOW);   
      delay(100);           
      digitalWrite(2, HIGH);    
      delay(100);  
  }
}

void loop() {
  if (!client.connected()) {
    reconnect();
  }
  client.loop();

  unsigned long now = millis();
  if (now - lastMsg > 2000) {
    lastMsg = now;
    client.publish(service_topic, clientId.c_str());
    blink(1);
  }
 
}

