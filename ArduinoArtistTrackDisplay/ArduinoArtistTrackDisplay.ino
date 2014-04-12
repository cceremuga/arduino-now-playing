#include <LiquidCrystal.h>

//initialize the library with the numbers of the interface pins
LiquidCrystal lcd(12, 11, 5, 4, 3, 2);

void setup() {
  //set up the LCD's number of columns and rows: 
  lcd.begin(16, 2);
  
  //initialize the serial communications:
  Serial.begin(9600);
  
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
    
    //wait a bit for the entire message to arrive
    delay(200);
    
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
      int newLineIndex = content.indexOf("<~>");
      
      String lineOne = content.substring(0, newLineIndex);
      String lineTwo = content.substring(newLineIndex + 3);
      
      lcd.print(lineOne);
      lcd.setCursor(0, 1);
      lcd.print(lineTwo);
      
      //wait 1 secs
      delay(1000);
      
      //scroll if either line > 16 chars
      if (lineOne.length() > 16 || lineTwo.length() > 16) {
        int scrollSize = 0;
        
        //figure out which line is longer and set that to the scroll size
        if(lineOne.length() > lineTwo.length()) {
          scrollSize = (lineOne.length() - 16);
        } else {
          scrollSize = (lineTwo.length() - 16);
        }
        
        //scroll 'em if you got 'em
        if(scrollSize > 0) {
          for(int i = 0;i < scrollSize; i++) {
            //scroll on back!
            lcd.scrollDisplayLeft(); 
            delay(200);
          }
          
          for(int i = 0;i < scrollSize; i++) {
            //scroll to the right
            lcd.scrollDisplayRight(); 
            delay(200);
          }
        }
      }
    }
  }
}
