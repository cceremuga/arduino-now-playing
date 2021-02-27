// Omit the Adafruit imports and adafruit LCD library
// specific code if just using standard LiquidCrysal.h
#include <Wire.h>
#include <Adafruit_RGBLCDShield.h>
//#include <LiquidCrystal.h>

#define BAUD_RATE 9600
#define DISPLAY_COLUMNS 16
#define DISPLAY_ROWS 2
#define SCROLL_DELAY 175
#define MESSAGE_WAIT 200
#define SCROLL_PAUSE 1000

Adafruit_RGBLCDShield lcd = Adafruit_RGBLCDShield();
// Initialize the library with the numbers of the interface pins
//LiquidCrystal lcd(12, 11, 5, 4, 3, 2);

void setup() {
  lcd.begin(DISPLAY_COLUMNS, DISPLAY_ROWS);
  Serial.begin(BAUD_RATE);

  lcd.setBacklight(0x2); //comment out this line if not using Adafruit RGB LCD
  lcd.noAutoscroll();
  lcd.print("Idle, waiting");
  lcd.setCursor(0, 1);
  lcd.print("for connection.");
  lcd.setCursor(0, 0);
}

void loop() {
  if (!Serial.available()) {
    return;
  }

  // Wait for the entire message to arrive,
  // may need to be adjusted if you notice odd output issues.
  delay(MESSAGE_WAIT);

  // Wipe the display.
  lcd.setCursor(0, 0);
  lcd.clear();

  String buffer = "";
  char character;

  // Fill the buffer.
  while (Serial.available() > 0) {
    character = Serial.read();
    buffer.concat(character);
  }

  if (buffer == "") {
    return;
  }

  int scrollSize = print(buffer);
  bounce(scrollSize);
}

void bounce(int scrollSize) {
  if (scrollSize == 0) {
    return;
  }

  while (Serial.available() == 0) {
    delay(SCROLL_PAUSE);
    scrollLeft(scrollSize);
    delay(SCROLL_PAUSE);
    scrollRight(scrollSize);
  }
}

int print(String serialContent) {
  // Split into two strings based upon ridiculous delimeter.
  int newLineIndex = serialContent.indexOf("<~>");

  String lineOne = serialContent.substring(0, newLineIndex);
  lcd.print(lineOne);

  String lineTwo = serialContent.substring(newLineIndex + 3);
  lcd.setCursor(0, 1);
  lcd.print(lineTwo);

  return scrollSize(lineOne, lineTwo);
}

int scrollSize(String lineOne, String lineTwo) {
  // Scroll directionally if either line > DISPLAY_COLUMNS.
  if (lineOne.length() > DISPLAY_COLUMNS || lineTwo.length() > DISPLAY_COLUMNS) {
    // Figure out which line is longer and set that to the scroll size.
    if (lineOne.length() > lineTwo.length()) {
      return (lineOne.length() - DISPLAY_COLUMNS);
    }

    return (lineTwo.length() - DISPLAY_COLUMNS);
  }

  return 0;
}

void scrollLeft(int numCharacters) {
  for (int i = 0; i < numCharacters; i++) {
    lcd.scrollDisplayLeft();
    delay(SCROLL_DELAY);
  }
}

void scrollRight(int numCharacters) {
  for (int i = 0; i < numCharacters; i++) {
    lcd.scrollDisplayRight();
    delay(SCROLL_DELAY);
  }
}
