#include <LiquidCrystal.h>

//initialize the library with the numbers of the interface pins
LiquidCrystal lcd(12, 11, 5, 4, 3, 2);

void setup(){
  //set up the LCD's number of columns and rows: 
  lcd.begin(16, 2);
  
  //initialize the serial communications:
  Serial.begin(9600);
  
  lcd.print("Serial active.");
  lcd.setCursor(0, 1);
  lcd.print("Waiting...");
  lcd.setCursor(0, 0);
}

void loop()
{  
  //when characters arrive over the serial port...
  if (Serial.available()) {
    //wait a bit for the entire message to arrive
    delay(100);
    
    //clear the screen
    lcd.noAutoscroll();
    lcd.setCursor(0, 0);
    lcd.clear();
    
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
      String lineTwo = content.substring(newLineIndex + 1);
      
      lcd.println(lineOne);
      lcd.print(lineTwo);
    }
    
    //read all the available characters
   // while (Serial.available() > 0) {
    //  char ch = Serial.read();
      
    //  if (ch == '\n')
    //  {
    //    lcd.setCursor(0, 1);
    //  }
    //  else
    //  {
        //display each character to the LCD
    //    lcd.write(ch);
    //  }
    //}
    
    delay(2000);
  }
}
