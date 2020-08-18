#include <CapacitiveSensor.h>
#include <Ewma.h>

#define FILTER 0.1

#define SAMPLES 5

#define DELAY 1 // 1KHz

CapacitiveSensor s1 = CapacitiveSensor(6, 0);
CapacitiveSensor s2 = CapacitiveSensor(6, 1);
CapacitiveSensor s3 = CapacitiveSensor(6, 2);
CapacitiveSensor s4 = CapacitiveSensor(6, 3);
CapacitiveSensor s5 = CapacitiveSensor(6, 4);
CapacitiveSensor s6 = CapacitiveSensor(6, 5);

Ewma f1(FILTER);
Ewma f2(FILTER);
Ewma f3(FILTER);
Ewma f4(FILTER);
Ewma f5(FILTER);
Ewma f6(FILTER);

void setup() {
  Serial.begin(115200);

  s1.set_CS_AutocaL_Millis(2000);
  s2.set_CS_AutocaL_Millis(2000);
  s3.set_CS_AutocaL_Millis(2000);
  s4.set_CS_AutocaL_Millis(2000);
  s5.set_CS_AutocaL_Millis(2000);
  s6.set_CS_AutocaL_Millis(2000);
}

void loop() {
  float r1 = (float) s1.capacitiveSensor(SAMPLES);
  float r2 = (float) s2.capacitiveSensor(SAMPLES);
  float r3 = (float) s3.capacitiveSensor(SAMPLES);
  float r4 = (float) s4.capacitiveSensor(SAMPLES);
  float r5 = (float) s5.capacitiveSensor(SAMPLES);
  float r6 = (float) s6.capacitiveSensor(SAMPLES);
  
  float v1 = f1.filter(r1);
  float v2 = f2.filter(r2);
  float v3 = f3.filter(r3);
  float v4 = f4.filter(r4);
  float v5 = f5.filter(r5);
  float v6 = f6.filter(r6);

  Serial.print(v1);
  Serial.print(",");
  Serial.print(v2);
  Serial.print(",");
  Serial.print(v3);
  Serial.print(",");
  Serial.print(v4);
  Serial.print(",");
  Serial.print(v5);
  Serial.print(",");
  Serial.print(v6);
  Serial.println();

  delay(DELAY);
}
