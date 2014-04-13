using log4net;
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

namespace NowPlayingToSerial
{
    class NowPlayingToSerial
    {
        private const String VLC_WEB_URL = "http://localhost:8080/requests/status.json";
        private const Int32 DEFAULT_VLC_POLL_RATE = 15000;
        private const Int32 DEFAULT_BAUD_RATE = 9600;
        private const String DEFAULT_SERIAL_PORT = "COM3";

        private Timer _vlcPollTimer;
        private String _customVlcUrl = String.Empty;
        private SerialPort _outputSerialPort;
        private String _lastSentMessage = String.Empty;
        private static readonly log4net.ILog log = log4net.LogManager.GetLogger (System.Reflection.MethodBase.GetCurrentMethod().DeclaringType);

        /// <summary>
        /// Constructor which grabs optional overrides via interactive input menu.
        /// </summary>
        public NowPlayingToSerial()
        {
            //banner
            DisplayBanner();

            //grab port from user
            String portName = ListAndPickASerialPort();

            //grab baud rate
            Int32 baudRate = PickABaudRate();

            //grab menu choice
            Int32 menuChoice = PresentMenuOptions();

            if (menuChoice == 1)
            {
                //grab VLC url from user
                _customVlcUrl = GetVlcUrlFromUser();
            }

            //serial time...
            if (baudRate > 0)
            {
                //open the serial port
                try
                {
                    _outputSerialPort = new SerialPort(portName, baudRate);
                    _outputSerialPort.Open();
                }
                catch (Exception ex)
                {
                    LogError(String.Format("{0} could not be connected to at {1} baud. Please run the program again to try again. Press enter to quit.", portName, baudRate), ex);
                }

                if (_outputSerialPort.IsOpen)
                {
                    Console.Clear();
                    log.Info(String.Format("Connected to {0} at {1} baud, here we go!", portName, baudRate));

                    if (menuChoice == 1)
                    {
                        try
                        {
                            //we're good to go!
                            InitializeVlc();
                        }
                        catch (WebException ex)
                        {
                            LogError(String.Format("Looks like VLC is unable to be found active at {0}. Please run the program again after starting VLC's web interface. Press enter to quit.", VLC_WEB_URL), ex);
                        }
                    }
                    else
                    {
                        Console.WriteLine("You've picked an unavailable option. Nothing left to do but press ENTER to quit.");
                    }
                }

            }
        }

        static void Main(string[] args)
        {
            //log4net setup
            log4net.Config.BasicConfigurator.Configure();
            ILog log = log4net.LogManager.GetLogger(typeof(NowPlayingToSerial));

            //CONSTRUCT
            NowPlayingToSerial p = new NowPlayingToSerial();

            //if they hit enter, let's quit.
            Console.ReadLine();

            //close down port just to be nice.
            p._outputSerialPort.Close();

            Environment.Exit(0);
        }

        /// <summary>
        /// Presents choices between VLC and Spotify
        /// </summary>
        /// <returns>Menu option</returns>
        private int PresentMenuOptions()
        {
            Int32 menuChoice = 1;

            Console.WriteLine("\nWhich would you like to track? (1) for VLC, or (2) for Spotify. [1]");

            String menuChoiceRaw = Console.ReadLine();

            //parse what they typed in...
            if (!String.IsNullOrEmpty(menuChoiceRaw))
            {
                Int32.TryParse(menuChoiceRaw, out menuChoice);
                Console.WriteLine();
            }

            return menuChoice == 0 ? 1 : menuChoice;
        }

        /// <summary>
        /// Overrides the default VLC URL with one specified by the user
        /// </summary>
        /// <returns>URL to the VLC web interface jSON page</returns>
        private string GetVlcUrlFromUser()
        {
            String vlcUrl = VLC_WEB_URL;

            Console.WriteLine(String.Format("\nWhat's the URL to your VLC web interface jSON page? [{0}]", VLC_WEB_URL));

            String overrideUrl = Console.ReadLine();

            //override vlc url if they put anything into the input.
            if (!String.IsNullOrEmpty(overrideUrl) && overrideUrl != VLC_WEB_URL)
                vlcUrl = overrideUrl;

            return vlcUrl;
        }

        /// <summary>
        /// Displays a silly ASCII banner.
        /// </summary>
        private void DisplayBanner()
        {
            Console.ForegroundColor = ConsoleColor.Cyan;
            Console.WriteLine("    _   __                 ____  __            _            ");
            Console.WriteLine("   / | / /___ _      __   / __ \\/ /___ ___  __(_)___  ____ _");
            Console.WriteLine("  /  |/ / __ \\ | /| / /  / /_/ / / __ `/ / / / / __ \\/ __ `/");
            Console.WriteLine(" / /|  / /_/ / |/ |/ /  / ____/ / /_/ / /_/ / / / / / /_/ / ");
            Console.WriteLine("/_/ |_/\\____/|__/|__/  /_/   /_/\\__,_/\\__, /_/_/ /_/\\__, /  ");
            Console.WriteLine("   __           _____           _    /____/        /____/   ");
            Console.WriteLine("  / /_____     / ___/___  _____(_)___ _/ /                  ");
            Console.WriteLine(" / __/ __ \\    \\__ \\/ _ \\/ ___/ / __ `/ /                   ");
            Console.WriteLine("/ /_/ /_/ /   ___/ /  __/ /  / / /_/ / /                    ");
            Console.WriteLine("\\__/\\____/   /____/\\___/_/  /_/\\__,_/_/                     \n\n");
            Console.ResetColor();
        }

        /// <summary>
        /// Grabs what's currently playing in VLC, sends it immediately, fires off a timer for timed retrievals from VLC
        /// </summary>
        private void InitializeVlc()
        {
            //fire off now!
            GrabNowPlayingFromVlc();

            //30 second timer
            _vlcPollTimer = new Timer(DEFAULT_VLC_POLL_RATE);
            _vlcPollTimer.Elapsed += new ElapsedEventHandler(VlcTimerElapsed);
            _vlcPollTimer.Enabled = true;
        }

        /// <summary>
        /// Interactively polls a user for the baud rate of their serial device
        /// </summary>
        /// <returns>Baud rate</returns>
        private Int32 PickABaudRate()
        {
            Int32 baudRate = DEFAULT_BAUD_RATE;

            Console.WriteLine(String.Format("\nWhat baud rate? [{0}]", DEFAULT_BAUD_RATE.ToString()));

            String baudRateRaw = Console.ReadLine();

            //parse baud rate if they changed from default
            if (!String.IsNullOrEmpty(baudRateRaw))
            {
                Int32.TryParse(baudRateRaw, out baudRate);
                Console.WriteLine();
            }

            return baudRate;
        }

        /// <summary>
        /// Interactively polls a user for the serial port name of their device. Lists available ports prior to gathering input
        /// </summary>
        /// <returns>Serial port name</returns>
        private String ListAndPickASerialPort()
        {
            String portName = DEFAULT_SERIAL_PORT;

            //get a list of serial port names. 
            string[] serialPorts = SerialPort.GetPortNames();

            Console.WriteLine("Found the following available serial port(s):\n");

            //display each serial port to the user
            foreach (string serialPort in serialPorts)
                Console.WriteLine(" " + serialPort);

            Console.WriteLine(String.Format("\nWhich one would you like to use? [{0}]", DEFAULT_SERIAL_PORT));

            //grab the serial port
            String userInputPortName = Console.ReadLine();

            if (!String.IsNullOrEmpty(userInputPortName) && userInputPortName != DEFAULT_SERIAL_PORT)
                portName = userInputPortName;

            return portName;
        }

        /// <summary>
        /// VLC timer elapsed event.
        /// </summary>
        /// <param name="sender"></param>
        /// <param name="e"></param>
        private void VlcTimerElapsed(object sender, ElapsedEventArgs e)
        {
            GrabNowPlayingFromVlc();
        }

        /// <summary>
        /// Grabs now playing jSON from VLC web interface, deserializes it, parses the now playing message
        /// only in the event that it is different from the previously sent message. This handles a use case
        /// where we're polling repeatedly while the same song is playing. No sense in sending the data twice.
        /// </summary>
        private void GrabNowPlayingFromVlc()
        {
            //grab jSON
            var nowPlayingJson = new WebClient().DownloadString(_customVlcUrl);

            //deserialize
            VlcStatus nowPlaying = JsonConvert.DeserializeObject<VlcStatus>(nowPlayingJson);

            //format
            String nowPlayingFormatted = nowPlaying.information.category.meta.now_playing.Replace(" - ", "<~>");

            //if different from last status, update last status, send to port
            if (_lastSentMessage != nowPlayingFormatted)
            {
                _lastSentMessage = nowPlayingFormatted;

                this.SendTextToSerial(_outputSerialPort, _lastSentMessage);
            }
        }

        /// <summary>
        /// Logs a given exception / outputs friendly error message to screen console.
        /// </summary>
        /// <param name="text">Friendly error message</param>
        /// <param name="ex">Raw exception</param>
        private void LogError(String text, Exception ex)
        {
            log.Error(ex.ToString());

            Console.ForegroundColor = ConsoleColor.Red;
            Console.WriteLine(text);
            Console.ResetColor();
        }

        /// <summary>
        /// Sends a string to the specified serial port.
        /// </summary>
        /// <param name="port">Instantiated / open serial port</param>
        /// <param name="text">Text to send</param>
        private void SendTextToSerial(SerialPort port, String text)
        {
            if (port != null && port.IsOpen)
            {
                //write to port
                port.Write(text);

                //log to console
                log.Info(String.Format("Sent: {0}", text));
            }
        }
    }
}
