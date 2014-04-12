using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System;
using System.Collections.Generic;
using System.IO.Ports;
using System.Linq;
using System.Net;
using System.Text;
using System.Threading.Tasks;
using System.Timers;

namespace VLCToSerial
{
    class Program
    {
        private Timer _timer;
        private SerialPort _port;
        private String VLC_WEB_URL = "http://localhost:8080/requests/status.json";
        private Int32 TIMER_MILLISECONDS = 30000;
        private Int32 DEFAULT_BAUD_RATE = 9600;
        private String _lastStatus = String.Empty;

        private SerialPort InitializeNowPlaying()
        {
            //get a list of serial port names. 
            string[] serialPorts = SerialPort.GetPortNames();

            Console.WriteLine("We found the following available serial ports:\n");

            //display each serial port to the user
            Console.ForegroundColor = ConsoleColor.DarkGreen;
            foreach (string serialPort in serialPorts)
                Console.WriteLine("  -" + serialPort);

            Console.ResetColor();
            Console.WriteLine("\nWhich one would you like to use?\n");

            //grab the serial port
            String portName = Console.ReadLine();

            //grab baud rate
            Console.WriteLine(String.Format("\nExcellent! What baud rate? [{0}]\n", DEFAULT_BAUD_RATE.ToString()));

            String baudRateRaw = Console.ReadLine();
            Int32 baudRate = DEFAULT_BAUD_RATE;

            //parse baud rate if they changed from default
            if (!String.IsNullOrEmpty(baudRateRaw))
            {
                Int32.TryParse(baudRateRaw, out baudRate);
                Console.WriteLine();
            }

            if (baudRate > 0)
            {
                //open the serial port
                _port = new SerialPort(portName, baudRate);
                _port.Open();

                if (_port.IsOpen)
                {
                    Console.Clear();
                    Console.WriteLine("Cool, we're connected to the serial port, here we go!\n");

                    //fire off now!
                    GrabNowPlayingFromVlc();

                    //30 second timer
                    _timer = new Timer(TIMER_MILLISECONDS);
                    _timer.Elapsed += new ElapsedEventHandler(timerElapsed);
                    _timer.Enabled = true;
                }
            }

            return _port;
        }

        private void timerElapsed(object sender, ElapsedEventArgs e)
        {
            GrabNowPlayingFromVlc();
        }

        private void GrabNowPlayingFromVlc()
        {
            //grab jSON
            var nowPlayingJson = new WebClient().DownloadString(VLC_WEB_URL);

            //deserialize
            VlcStatus nowPlaying = JsonConvert.DeserializeObject<VlcStatus>(nowPlayingJson);

            //format
            String nowPlayingFormatted = nowPlaying.information.category.meta.now_playing.Replace(" - ", "<~>");

            //if different from last status, update last status, send to port
            if (_lastStatus != nowPlayingFormatted)
            {
                _lastStatus = nowPlayingFormatted;

                this.SendTextToSerial(_port, _lastStatus);
            }
            else
            {
                this.LogToConsole(String.Format("Nothing new to send, sleeping for {0} milliseconds.", TIMER_MILLISECONDS));
            }
        }

        private void LogToConsole(String text)
        {
            Console.WriteLine(String.Format("[{0}] {1}", DateTime.Now, text));
        }

        private void SendTextToSerial(SerialPort port, String text)
        {
            if (port != null && port.IsOpen)
            {
                //write to port
                port.Write(text);

                //log to console
                this.LogToConsole(String.Format("Sent: {0}", text));
            }
        }

        static void Main(string[] args)
        {
            Program p = new Program();

            //fire it up
            p.InitializeNowPlaying();

            Console.ReadLine();

            p._port.Close();
        }
    }
}
