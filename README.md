![Arduino Now Playing](http://i.imgur.com/NuFMnNg.png "Arduino Now Playing")
=========

Arduino Now Playing is a small project to provide source code to developers, enabling them to send "Now Playing" audio player metadata to an [Arduino] connected 16x2 character LCD display.

Included is also a sample client written in Go for periodically polling play status from an authorized Spotify account.

On the receiving end, an Arduino sketch listens via serial port, receives data, splits artist / track into two separate lines and displays via `LiquidCrystal.h`. Scrolling, if the length exceeds your usable display width.

Release Notes
----

* **2.26.2021** - Refactoring, cleanup, upgrade to Go 1.16, Spotify support.
* **4.26.2014** - Go port stable.

Spotify Configuration
----

1. Create a Spotify developer account, app, retrieve OAuth client id, secret.
2. Set environment variable values: `SPOTIFY_ID`, `SPOTIFY_SECRET`.
3. On your Spotify developer app, ensure the callback URL is whitelisted: `http://localhost:8080/callback`.
4. Run the Go spotify-client.

Roadmap
----

* Unit tests, probably.
* Additional support in Arduino sketch for very long titles.
* The sky's the limit, share your suggestions, please!

Flashy action shots
----

A prototype running on a [Sparkfun RedBoard] receiving serial data and displaying on a 16x2 LCD character display:

![Basic Prototype](http://i.imgur.com/cSCjJos.jpg "Basic Prototype")

My completed build featuring an [Adafruit display] running on an Arduino Uno 

![Completed Build](http://i.imgur.com/jw8FG55.jpg "Completed Build")

License
-----------

MIT. See `LICENSE.md` for more info.

[Arduino]:http://arduino.cc/
[Sparkfun RedBoard]:https://www.sparkfun.com/products/11575
[Adafruit display]:http://www.adafruit.com/products/716