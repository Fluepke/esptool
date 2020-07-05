package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const version string = "0.1"
const defaultConnectBaudrate uint = 115200
const defaultTransferBaudrate uint = 921600

type CliCommand struct {
	Name        string
	Description string
	FlagSet     *flag.FlagSet
	Callback    func(*log.Logger) error
}

var (
	versionFlagSet = flag.NewFlagSet("version", flag.ExitOnError)
	versionJson    = versionFlagSet.Bool("json", false, "Display version info in JSON format")

	help = flag.Bool("help", false, "Show a help page")

	infoFlagSet          = flag.NewFlagSet("info", flag.ExitOnError)
	infoPort             = infoFlagSet.String("serial.port", "", "Serial port device file")
	infoConnectBaudrate  = infoFlagSet.Uint("serial.baudrate.connect", defaultConnectBaudrate, "Serial signalling rate during connect phase")
	infoTransferBaudrate = infoFlagSet.Uint("serial.baudrate.transfer", defaultTransferBaudrate, "Serial signalling rate during data transfer")
	infoTimeout          = infoFlagSet.Duration("serial.connect.timeout", 500*time.Millisecond, "Timeout to wait for chip response upon connecting")
	infoRetries          = infoFlagSet.Uint("serial.connect.retries", 5, "How often to retry connecting")
	infoJson             = infoFlagSet.Bool("json", false, "Display chip info in JSON format")

	flashReadFlagSet          = flag.NewFlagSet("readFlash", flag.ExitOnError)
	flashReadPort             = flashReadFlagSet.String("serial.port", "", "Serial port device file")
	flashReadConnectBaudrate  = flashReadFlagSet.Uint("serial.baudrate.connect", defaultConnectBaudrate, "Serial signalling rate during connect phase")
	flashReadTransferBaudrate = flashReadFlagSet.Uint("serial.baudrate.transfer", defaultTransferBaudrate, "Serial signalling rate during data transfer")
	flashReadTimeout          = flashReadFlagSet.Duration("serial.connect.timeout", 500*time.Millisecond, "Timeout to wait for chip response upon connecting")
	flashReadRetries          = flashReadFlagSet.Uint("serial.connect.retries", 5, "How often to retry connecting")
	flashReadOffset           = flashReadFlagSet.Uint("flash.offset", 0, "Offset")
	flashReadSize             = flashReadFlagSet.Uint("flash.size", 0, "Bytes to read")
	flashReadFile             = flashReadFlagSet.String("flash.file", "", "File to read flash contents into")
	flashReadPartitionName    = flashReadFlagSet.String("flash.partition.name", "", "Partition to read")

	flashWriteFlagSet          = flag.NewFlagSet("writeFlash", flag.ExitOnError)
	flashWritePort             = flashWriteFlagSet.String("serial.port", "", "Serial port device file")
	flashWriteConnectBaudrate  = flashWriteFlagSet.Uint("serial.baudrate.connect", defaultConnectBaudrate, "Serial signalling rate during connect phase")
	flashWriteTransferBaudrate = flashWriteFlagSet.Uint("serial.baudrate.transfer", defaultTransferBaudrate, "Serial signalling rate during data transfer")
	flashWriteTimeout          = flashWriteFlagSet.Duration("serial.connect.timeout", 500*time.Millisecond, "Timeout to wait for chip response upon connecting")
	flashWriteRetries          = flashWriteFlagSet.Uint("serial.connect.retries", 5, "How often to retry connecting")
	flashWriteOffset           = flashWriteFlagSet.Uint("flash.offset", 0, "Offset")
	flashWriteFile             = flashWriteFlagSet.String("flash.file", "", "File with data to flash")
	flashWritePartitionName    = flashWriteFlagSet.String("flash.partition.name", "", "Partition to write")
	flashWriteCompress         = flashWriteFlagSet.Bool("flash.compress", true, "Use compression for transfer")

	cliCommands = []*CliCommand{
		&CliCommand{
			Name:        "version",
			Description: "Show version info and exit",
			FlagSet:     versionFlagSet,
			Callback: func(logger *log.Logger) error {
				versionFlagSet.Parse(os.Args[2:])
				return versionCommand(*versionJson)
			},
		},
		&CliCommand{
			Name:        "info",
			Description: "Retrieve various information from chip",
			FlagSet:     infoFlagSet,
			Callback: func(logger *log.Logger) error {
				infoFlagSet.Parse(os.Args[2:])
				esp32, err := connectEsp32(*infoPort, uint32(*infoConnectBaudrate), uint32(*infoTransferBaudrate), *infoRetries, logger)
				if err != nil {
					return err
				}
				return infoCommand(*infoJson, esp32)
			},
		},
		&CliCommand{
			Name:        "flashRead",
			Description: "Read flash contents",
			FlagSet:     flashReadFlagSet,
			Callback: func(logger *log.Logger) error {
				flashReadFlagSet.Parse(os.Args[2:])
				esp32, err := connectEsp32(*flashReadPort, uint32(*flashReadConnectBaudrate), uint32(*flashReadTransferBaudrate), *flashReadRetries, logger)
				if err != nil {
					return err
				}
				bytes, err := esp32.ReadFlash(uint32(*flashReadOffset), uint32(*flashReadSize))
				if err != nil {
					return err
				}
				os.Stdout.Write(bytes)
				return nil
			},
		},
		&CliCommand{
			Name:        "flashWrite",
			Description: "Write flash contents",
			FlagSet:     flashWriteFlagSet,
			Callback: func(logger *log.Logger) error {
				flashWriteFlagSet.Parse(os.Args[2:])
				contents, err := ioutil.ReadFile(*flashWriteFile)
				if err != nil {
					return err
				}
				esp32, err := connectEsp32(*flashWritePort, uint32(*flashWriteConnectBaudrate), uint32(*flashWriteTransferBaudrate), *flashWriteRetries, logger)
				if err != nil {
					return err
				}

				err = esp32.WriteFlash(uint32(*flashWriteOffset), contents, *flashWriteCompress)
				if err != nil {
					panic(err)
				}
				logger.Print("Done")
				return nil
			},
		},
	}

	port            = flag.String("port", "", "Serial port device")
	baudrate        = flag.Uint("baudrate", 115200, "Serial port baud rate used when flashing/reading")
	connectTimeout  = flag.Duration("timeout.connect", 500*time.Millisecond, "Timeout to wait for chip response upon connecting")
	responseTimeout = flag.Duration("timeout.response", 10*time.Millisecond, "Timeout to wait for chip to respond")

	server        = flag.NewFlagSet("server", flag.ExitOnError)
	listenAddress = server.String("listen-address", ":8080", "Port and andress to listen on")
)

func main() {
	flag.Parse()
	if len(os.Args) < 2 || *help {
		printHelp()
	}

	logger := log.New(os.Stderr, bold("[LOG]: "), log.Ltime|log.Lshortfile)

	for _, command := range cliCommands {
		if command.Name == os.Args[1] {
			err := command.Callback(logger)
			if err != nil {
				logger.Printf("Failed to run %s: %s", command.Name, err.Error())
			}
			os.Exit(1)
		}
	}
}

func printHelp() {
	fmt.Println("Please choose one of the following subcommands")
	for _, command := range cliCommands {
		fmt.Printf("  * \033[1m%s\033[0m: %s\n", command.Name, command.Description)
	}
	os.Exit(1)
}
