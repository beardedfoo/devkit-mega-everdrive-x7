package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	serial "github.com/goburrow/serial"
)

// Parameters from Mega Everdrive X7 design
const blockSize = 512 * 128
const maxGameSize = 0xf00000

// Syntax specific to the Mega Everdrive X7 cartridge
const x7_InitCmd = "    *T"
const x7_LoadGameCmd = "*g"
const x7_OK = "k"
const x7_DataOK = "d"
const x7_RunSMS = "*rs"
const x7_RunOS = "*ro"
const x7_RunCD = "*rc"
const x7_RunM10 = "*rM"
const x7_RunSSF = "*rS"
const x7_RunMegadrive = "*rm"

// Flags from command line
var baudRate = 0
var serialPort = ""
var readTimeout = 0
var runMode = ""

func init() {
	flag.IntVar(&baudRate, "baudRate", 9600, "Serial baud rate")
	flag.IntVar(&readTimeout, "readTimeout", 1000, "Serial read timeout in msec")
	flag.StringVar(&serialPort, "serialPort", "/dev/tty.usbserial-A50543G8", "Serial port for Mega Everdrive X7")
	flag.StringVar(&runMode, "runMode", "md", "Sets the run mode for the rom: sms|cd|os|md|m10|ssf")
}

func main() {
	var gameData []byte
	var err error

	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf("Usage: megaed-run [options] rom\n")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	// Get the rom file to load as a positional argument
	romFile := flag.Arg(0)

	// Read the rom data from disk
	gameData, err = ioutil.ReadFile(romFile)
	if err != nil {
		log.Fatalf("Loading game data failed: %v", err)
		os.Exit(-1)
	} else {
		fmt.Printf("Read %d bytes from rom file\n", len(gameData))

		// Ensure the game data ends with a full block
		padding := 0
		for len(gameData)%blockSize != 0 {
			padding += 1
			gameData = append(gameData, 0x00)
		}

		if len(gameData) > maxGameSize {
			log.Fatalf("Game data exceeds maximum size of %d bytes", maxGameSize)
			os.Exit(-1)
		}
	}

	// Check the ROM
	if string(gameData[0x100:0x104]) != "SEGA" {
		fmt.Printf("WARNING: ROM may be corrupt: expected string 'SEGA' at offset 0x100\n")
	}

	// Open serial port (optional)
	fmt.Printf("Connecting to serial...")
	serialConfig := &serial.Config{
		Address: serialPort, BaudRate: baudRate, StopBits: 1,
		Timeout: 1000 * time.Millisecond, Parity: "N"}
	fSerial, err := serial.Open(serialConfig)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
		os.Exit(1)
	} else {
		fmt.Printf("OK\n")
	}
	defer fSerial.Close()

	// Check that the cart communication is working
	checkX7(fSerial)

	// Send the game data to the console
	sendGame(fSerial, gameData)

	// Send the command to start the game
	startGame(fSerial, runMode)
}

func startGame(fh io.ReadWriter, runMode string) {
	fmt.Printf("Starting game...")
	// Send the selected game start command
	switch runMode {
	case "md":
		writeSerial(fh, []byte(x7_RunMegadrive))
	case "sms":
		writeSerial(fh, []byte(x7_RunSMS))
	case "m10":
		writeSerial(fh, []byte(x7_RunM10))
	case "cd":
		writeSerial(fh, []byte(x7_RunCD))
	case "os":
		writeSerial(fh, []byte(x7_RunOS))
	case "ssf":
		writeSerial(fh, []byte(x7_RunSSF))
	default:
		fmt.Printf("Unsupported runMode specified\n")
		os.Exit(-1)
	}

	// Receive the response
	if resp := readSerial(fh); resp != x7_OK {
		fmt.Printf("ERROR\n")
		log.Fatalf("ERROR: Bad response from cart: %v", resp)
		os.Exit(-1)
	}
	fmt.Printf("OK\n")
}

func sendGame(fh io.ReadWriter, gameData []byte) {
	expectedMd5 := md5.New()
	expectedMd5.Write(gameData)

	// Send the command to load game data followed by the number of expected blocks
	writeSerial(fh, []byte(x7_LoadGameCmd))
	blockCount := len(gameData) / blockSize
	txData := []byte{byte(blockCount)}
	writeSerial(fh, txData)

	// Receive the response
	if readSerial(fh) != x7_OK {
		log.Fatalf("ERROR: Bad response from cart")
		os.Exit(-1)
	}

	// Send the game data
	fmt.Printf("Sending game data")
	wroteMd5 := md5.New()
	for x := 0; x < blockCount; x++ {
		// Load the block as a byte buffer
		block := gameData[:blockSize]
		gameData = gameData[blockSize:]
		wroteMd5.Write(block)

		// Write the block to the tty
		writeSerial(fh, block)

		// Log progress to the user
		fmt.Printf(".")
	}

	// Check the md5 of the written data
	if !bytes.Equal(expectedMd5.Sum(nil), wroteMd5.Sum(nil)) {
		log.Printf("ERROR: bad md5: %x != %x", expectedMd5.Sum(nil), wroteMd5.Sum(nil))
	}

	// Confirm the console received the data block
	if readSerial(fh) != x7_DataOK {
		log.Fatalf("ERROR: Bad response from cart")
		os.Exit(-1)
	}

	fmt.Printf("OK\n")
}

// Verify communication with the Mega Everdrive X7
func checkX7(fh io.ReadWriter) {
	// Send the test command
	fmt.Printf("Connection test...")
	writeSerial(fh, []byte(x7_InitCmd))

	// Receive the response
	if readSerial(fh) != x7_OK {
		fmt.Printf("ERROR\n")
		log.Fatalf("ERROR: Bad response from cart")
		os.Exit(-1)
	}
	fmt.Printf("OK\n")
}

func writeSerial(fh io.Writer, buf []byte) {
	wrote := 0
	for wrote < len(buf) {
		n, err := fh.Write(buf[wrote:])
		wrote += n

		if err != nil {
			// Back off when needed
			if err.Error() == "resource temporarily unavailable" {
				time.Sleep(100 * time.Millisecond)
				log.Printf("trying again: %v", err)
				continue
			} else {
				log.Fatalf("write failed: %v", err)
				os.Exit(-1)
			}
		}

		// Sleep to avoid overloading the cartridge with data and generating
		// "resource temporarily unavailable" errors
		time.Sleep(1 * time.Millisecond)
	}
}

func readSerial(fh io.Reader) string {
	scanner := bufio.NewScanner(fh)
	scanner.Scan()
	return scanner.Text()
}
