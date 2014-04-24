What's this Arduino Now Playing thing?
=========

Arduino Now Playing is a small project whose goal is to provide source code to developers, enabling them to send "Now Playing" audio player updates to an [Arduino] connected 16x2 character LCD display.

A handy example Go app runs, grabs metadata, parses it and pushes it over the serial port of your choosing. Currently supported players are the fantastic [VLC] player via their jSON Web Interface and [Spotify]. The VLC integration was written with [soma fm] in mind, but would likely apply to other streams / local file play too. All you'd need to do is fork the repo and make your updates!

Included is a sample Arduino C++ sketch which listens via serial port, receives data, splits artist / track into two separate lines and displays via LiquidCrystal.h. Scrolling, if the length exceeds your usable display width.

Ultimitely, this is just a little project which I plan on taking from prototype to completion for personal use, but I figured others may find it interesting as well.

The newest features
----

* 4.23.2013 - Start of C# app port to Go. Henceforth the C# code is now deprecated in favor of a more cross platform friendly language. Default branch of repo has been switched to go-port. **For a more stable version, please opt for Master until dev / merge is complete.**
* 4.22.2014 - Spotify support added to C# example console app! Required library: [https://code.google.com/p/spotify-local-api]

Future improvements
----

* Additional support in Arduino C++ sketch for very long titles.
* The sky's the limit, share your suggestions, please!

Flashy action shots
----

Example C# console app displaying a running log of sent serial data:

![C# Console App](http://i.imgur.com/EKAqgqH.jpg "C# Console App")

A prototype Arduino Uno compatible board receiving serial data and displaying on a 16x2 LCD character display:

![LCD Character Display](http://i.imgur.com/cSCjJos.jpg "LCD Character Display")

Legal stuff
-----------

Copyright (c) 2014 Craig Ceremuga

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

[VLC]:http://www.videolan.org/vlc/index.html
[soma fm]:http://somafm.com/
[Arduino]:http://arduino.cc/
[https://code.google.com/p/spotify-local-api]:https://code.google.com/p/spotify-local-api
[Spotify]:https://www.spotify.com/