VLC -> Serial -> Arduino
=========

This project hooks into the [VLC] jSON Web Interface, grabs now_playing data, parses it, pushes it over the serial port of your choosing.

This was written with [soma fm] in mind, but would likely apply to other streams / local file play. All you'd need to do is fork the repo, make your updates & profit!

Included is a sample [Arduino] project which listens, receives serial data, splits artist / track into two separate lines and displays on a standard [16x2 character LED display].

Ultimitely, this is just a little project which I plan on taking from prototype to completion for personal use, but I figured others may find it interesting as well.

What's in store down the road?
----

  - Fix occasional timing issues with serial communication.
  - Less hard coded sillyness.
  - Actual versioning.
  - More interactive options in C# portion on startup.
  - Logging to file to keep tabs on what was previously played.
  - A project rename and potential Spotify support.
  - The sky's the limit, share your suggestions, please!

Flashy Action Shot
----

![LCD Output](http://i.imgur.com/cSCjJos.jpg "LCD Output")

License
-----------

The MIT License (MIT)

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

[vlc]:http://www.videolan.org/vlc/index.html
[soma fm]:http://somafm.com/
[Arduino]:http://arduino.cc/
[16x2 character LCD display]:https://www.sparkfun.com/products/709