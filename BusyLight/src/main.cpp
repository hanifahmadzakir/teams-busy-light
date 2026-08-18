#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>

#define SCREEN_WIDTH 128
#define SCREEN_HEIGHT 64
#define OLED_RESET -1
Adafruit_SSD1306 display(SCREEN_WIDTH, SCREEN_HEIGHT, &Wire, OLED_RESET);

// LED Pins (Wemos D1 mapping)
const int pinGreen = D5;
const int pinYellow = D6;
const int pinRed = D7;

// --- Function Prototypes for PlatformIO ---
void processCommand(String cmd);
void turnOnLed(int activePin);
void turnOffAllLeds();
void updateDisplay(String text);

// Built-in LED for visual feedback (Active LOW on Wemos D1)
const int pinBuiltIn = LED_BUILTIN; 

void setup() {
  Serial.begin(9600); // 9600 baud to match the Go backend
  
  pinMode(pinGreen, OUTPUT);
  pinMode(pinYellow, OUTPUT);
  pinMode(pinRed, OUTPUT);
  pinMode(pinBuiltIn, OUTPUT);
  
  // Keep built-in LED off (Active LOW means HIGH is off)
  digitalWrite(pinBuiltIn, HIGH);
  
  // Initialize OLED
  if(!display.begin(SSD1306_SWITCHCAPVCC, 0x3C)) {
    Serial.println("ERROR: SSD1306 OLED allocation failed");
    for(;;);
  }
  
  display.clearDisplay();
  display.setTextSize(2);
  display.setTextColor(SSD1306_WHITE);
  updateDisplay("Starting Up");
  
  // Tell the Wails App that the board is ready
  Serial.println("SYSTEM: Wemos D1 Initialized and Ready.");
  
  // Test LEDs on boot
  turnOnLed(pinGreen); delay(300);
  turnOnLed(pinYellow); delay(300);
  turnOnLed(pinRed); delay(300);
  turnOffAllLeds();
  
  updateDisplay("Waiting...");
}

void loop() {
  // Check if Wails app sent a command over USB Serial
  if (Serial.available() > 0) {
    String command = Serial.readStringUntil('\n');
    command.trim(); // Remove any whitespace/carriage returns
    
    if (command.length() > 0) {
      // Flash built-in LED to indicate data was received
      digitalWrite(pinBuiltIn, LOW); 
      
      processCommand(command);
      
      delay(50); // Short delay for visible blink
      digitalWrite(pinBuiltIn, HIGH);
    }
  }
}

void processCommand(String cmd) {
  if (cmd == "Green") {
    turnOnLed(pinGreen);
    updateDisplay("Available");
    Serial.println("DEBUG: Executed [Green] -> Available");
  } 
  else if (cmd == "Red") {
    turnOnLed(pinRed);
    updateDisplay("Busy / DND");
    Serial.println("DEBUG: Executed [Red] -> Busy / DND");
  } 
  else if (cmd == "Yellow") {
    turnOnLed(pinYellow);
    updateDisplay("Away");
    Serial.println("DEBUG: Executed [Yellow] -> Away");
  } 
  else if (cmd == "BlinkRed") {
    Serial.println("DEBUG: Executed [BlinkRed] -> Incoming Call!");
    updateDisplay("CALLING!");
    
    // Blink BOTH the external Red LED and the Built-in LED 5 times
    for(int i = 0; i < 5; i++) {
      turnOnLed(pinRed); 
      digitalWrite(pinBuiltIn, LOW); // Turn ON built-in LED (Active LOW)
      delay(200);
      
      turnOffAllLeds(); 
      digitalWrite(pinBuiltIn, HIGH); // Turn OFF built-in LED
      delay(200);
    }
    
    turnOnLed(pinRed);              // Return external LED to solid red
    digitalWrite(pinBuiltIn, HIGH); // Ensure built-in LED stays OFF after the call rings
  }
  else if (cmd == "Off") {
    turnOffAllLeds();
    updateDisplay("");
    Serial.println("DEBUG: Executed [Off] -> All LEDs off");
  }
  else if (cmd == "White") {
    turnOnLed(pinGreen); 
    updateDisplay("Connected!");
    Serial.println("DEBUG: Executed [White] -> Connection Established");
  }
  else {
    Serial.println("DEBUG: Unknown command received: " + cmd);
  }
}

void turnOnLed(int activePin) {
  turnOffAllLeds();
  digitalWrite(activePin, HIGH);
}

void turnOffAllLeds() {
  digitalWrite(pinGreen, LOW);
  digitalWrite(pinYellow, LOW);
  digitalWrite(pinRed, LOW);
}

void updateDisplay(String text) {
  display.clearDisplay();
  display.setCursor(0, 20);
  
  // Center text roughly depending on length (approx 12 pixels per char at size 2)
  int cursorX = (SCREEN_WIDTH - (text.length() * 12)) / 2;
  if (cursorX < 0) cursorX = 0;
  
  display.setCursor(cursorX, 24);
  display.println(text);
  display.display();
}