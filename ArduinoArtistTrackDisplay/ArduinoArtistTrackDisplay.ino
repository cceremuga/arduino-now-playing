#include <LiquidCrystal.h>

//configurations
#define BAUD_RATE 9600
#define DISPLAY_COLUMNS 16
#define DISPLAY_ROWS 2
#define SCROLL_DELAY 175
#define DELAY_BEFORE_SCROLL 1000
#define MESSAGE_WAIT 200

//initialize the library with the numbers of the interface pins
LiquidCrystal lcd(12, 11, 5, 4, 3, 2);

void setup() {
  //set up the LCD's number of columns and rows: 
  lcd.begin(DISPLAY_COLUMNS, DISPLAY_ROWS);
  
  //initialize the serial communications:
  Serial.begin(BAUD_RATE);
  
  //print something once we're ready to roll just to keep the user informed
  lcd.noAutoscroll();
  lcd.print("Serial active.");
  lcd.setCursor(0, 1);
  lcd.print("Waiting...");
  lcd.setCursor(0, 0);
}

void loop() {  

  //when characters arrive over the serial port...
  if (Serial.available()) {

    //wait for the entire message to arrive
    delay(MESSAGE_WAIT);
    
    //clear the screen
    lcd.setCursor(0, 0);
    lcd.clear();
    
    //this will act as our buffer as we read from serial
    String content = "";
    char character;
    
    //read all available into a string buffer
    while(Serial.available() > 0) {
      character = Serial.read();
      content.concat(character);
    }
    
    if (content != "") {
      //split into two strings
      int newLineIndex = content.indexOf("<~>");
      
      String lineOne = content.substring(0, newLineIndex);
      String lineTwo = content.substring(newLineIndex + 3);
      
      //print to display
      lcd.print(lineOne);
      lcd.setCursor(0, 1);
      lcd.print(lineTwo);
      
      //scroll directionally if either line > DISPLAY_COLUMNS
      if (lineOne.length() > DISPLAY_COLUMNS || lineTwo.length() > DISPLAY_COLUMNS) {
        int scrollSize = 0;
        
        //figure out which line is longer and set that to the scroll size
        if(lineOne.length() > lineTwo.length()) {
          scrollSize = (lineOne.length() - DISPLAY_COLUMNS);
        }
        else {
          scrollSize = (lineTwo.length() - DISPLAY_COLUMNS);
        }
        
        //scroll 'em if you got 'em
        if(scrollSize > 0) {
          //wait before scrolling.
          delay(DELAY_BEFORE_SCROLL);

          //now scroll left
          scrollLeft(scrollSize);

          //then back right
          scrollRight(scrollSize);
        }
      }
    }
  }
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