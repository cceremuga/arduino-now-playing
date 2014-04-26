//omit the adafruit imports and adafruit LCD library 
//specific code if just using standard LiquidCrysal.h
#include <Wire.h>
#include <Adafruit_MCP23017.h>
#include <Adafruit_RGBLCDShield.h>
//exclude above and include this if not using an 
//Adafruit RGB LCD
//#include <LiquidCrystal.h>

// The shield uses the I2C SCL and SDA pins. On classic Arduinos
// this is Analog 4 and 5 so you can't use those for analogRead() anymore
// However, you can connect other I2C sensors to the I2C bus and share
// the I2C bus.

//configuration
#define BAUD_RATE 9600
#define DISPLAY_COLUMNS 16
#define DISPLAY_ROWS 2
#define SCROLL_DELAY 175
#define DELAY_BEFORE_SCROLL 1000
#define MESSAGE_WAIT 200
#define SCROLL_PAUSE 1000
#define RED 0x1
#define YELLOW 0x3
#define GREEN 0x2
#define TEAL 0x6
#define BLUE 0x4
#define VIOLET 0x5
#define WHITE 0x7

//comment this out and use the below for non-Adafruit displays
Adafruit_RGBLCDShield lcd = Adafruit_RGBLCDShield();
//initialize the library with the numbers of the interface pins
//LiquidCrystal lcd(12, 11, 5, 4, 3, 2);

void setup() {
  //set up the LCD's number of columns and rows: 
  lcd.begin(DISPLAY_COLUMNS, DISPLAY_ROWS);
  
  //initialize the serial communications:
  Serial.begin(BAUD_RATE);
  
  //print something once we're ready to roll just to keep the user informed
  lcd.setBacklight(BLUE); //comment out this line if not using Adafruit RGB LCD
  lcd.noAutoscroll();
  lcd.print("Serial active.");
  lcd.setCursor(0, 1);
  lcd.print("Waiting...");
  lcd.setCursor(0, 0);
}

void loop() {  

  //when characters arrive over the serial port...
  if (Serial.available()) {

    //wait for the entire message to arrive, 
    //this timing may need to be adjusted if you notice odd output issues
    delay(MESSAGE_WAIT);
    
    //clear the screen to make way for the new message
    lcd.setCursor(0, 0);
    lcd.clear();
    
    //this will act as our buffer as we read from serial
    String serialBuffer = "";
    char character;
    
    //read all available into a string buffer
    while(Serial.available() > 0) {
      character = Serial.read();
      serialBuffer.concat(character);
    }
    
    //verify we've got a completed message
    if (serialBuffer != "") {
      //print to display
      int scrollSize = printToDisplay(serialBuffer);

      //scroll 'em if you got 'em
      if(scrollSize > 0) {
        //wait before scrolling just to give users a bit
        //to consume what has just appeared
        delay(DELAY_BEFORE_SCROLL);

        //scroll left
        scrollLeft(scrollSize);

        //wait to let the user consume before scrolling back
        delay(SCROLL_PAUSE);

        //then back right
        scrollRight(scrollSize);
      }
    }
  }
}

int printToDisplay(String serialContent) {
  //split into two strings based upon ridiculous delimeter
  int newLineIndex = serialContent.indexOf("<~>");

  String lineOne = serialContent.substring(0, newLineIndex);
  String lineTwo = serialContent.substring(newLineIndex + 3);

  //print to display, moving to the second line for line two
  lcd.print(lineOne);
  lcd.setCursor(0, 1);
  lcd.print(lineTwo);

  //determine number of characters to scroll, if any
  return determineScrollCharacters(lineOne, lineTwo);
}

int determineScrollCharacters(String lineOne, String lineTwo) {
  int scrollSize = 0;

  //scroll directionally if either line > DISPLAY_COLUMNS
  if (lineOne.length() > DISPLAY_COLUMNS || lineTwo.length() > DISPLAY_COLUMNS) {

    //figure out which line is longer and set that to the scroll size
    if(lineOne.length() > lineTwo.length()) {
      scrollSize = (lineOne.length() - DISPLAY_COLUMNS);
    }
    else {
      scrollSize = (lineTwo.length() - DISPLAY_COLUMNS);
    }
  }

  return scrollSize;
}

void scrollLeft(int numCharacters) {
  for(int i = 0;i < numCharacters; i++) {
    //scroll on back!
    lcd.scrollDisplayLeft(); 

    delay(SCROLL_DELAY);
  }
}

void scrollRight(int numCharacters) {
  for(int i = 0;i < numCharacters; i++) {
    //scroll to the right
    lcd.scrollDisplayRight(); 

    delay(SCROLL_DELAY);
  }
}
