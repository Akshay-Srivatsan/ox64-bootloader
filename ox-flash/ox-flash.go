// Utility for automatically flashing firmware onto the VisionFive 2

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jacobsa/go-serial/serial"

	_ "embed"
)

//go:embed clock_para.bin
var clock_para []byte

//go:embed flash_para.bin
var flash_para []byte

func exists(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

var baud = flag.Uint("baud", 2000000, "baud rate")
var file = flag.String("port", "", "serial port")
var run = flag.Bool("run", false, "exit the bootrom and run the program from flash")
var echo = flag.Bool("echo", false, "echo bytes received from the serial port after flashing")
var echo_baud = flag.Uint("echo-baud", 115200, "baud rate for echoing")
var bootloader = flag.Bool("bootloader", false, "flash a bootloader")
var payload = flag.Bool("payload", false, "flash a payload")
var addr = flag.Uint("addr", 0, "the address to flash to (default: 0x0 for firmware, 0x10000 for a payload)")

var autodetect = []string{
	"/dev/ttyACM0",
	"/dev/ttyUSB0",
	"/dev/ttyUSB1",
}

func main() {
	flag.Parse()

	if *file == "" {
		for _, f := range autodetect {
			if exists(f) {
				file = &f
				break
			}
		}
		if *file == "" {
			log.Fatal("could not autodetect serial port")
		}
	}

	if !*bootloader && !*payload {
		log.Fatal("must specify one of --bootloader or --payload")
	}

	if *echo && !*run {
		log.Fatal("the --echo option requires the --run option")
	}

	filetype := "firmware"
	if *payload {
		filetype = "payload"
	}

	if *addr == 0 {
		if *payload {
			*addr = 0x10000
		} else {
			*addr = 0x0
		}
	}

	args := flag.Args()
	if len(args) <= 0 {
		log.Fatalf("no %s file given", filetype)
	}
	if len(args) >= 2 {
		log.Fatalf("too many %s files given", filetype)
	}
	fw := args[0]

	options := serial.OpenOptions{
		PortName:        *file,
		BaudRate:        *baud,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("uart open: %v", err)
	}

	log.Printf("Connected to %s, baud: %d\n", options.PortName, options.BaudRate)

	jedec, err := handshake(port, options.BaudRate)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Have %s file: %s\n", filetype, fw)
	fw_data, err := os.ReadFile(fw)
	if err != nil {
		log.Fatal(err)
	}

	if !*payload {
		info, err := CreateBootInfo(jedec, fw_data, uint32(*addr))
		if err != nil {
			log.Fatal(err)
		}

		bh := bytes.NewBuffer(make([]byte, 0, 4096))
		err = binary.Write(bh, binary.LittleEndian, &info)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Writing bootrom-compatible header at 0x0000")
		err = flash(port, 0, bh.Bytes())
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Writing firmware at 0x%x\n", *addr)
		err = flash(port, uint32(*addr)+0x2000, fw_data)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("Writing payload at 0x%x\n", *addr)
		err = flash(port, uint32(*addr)+0x2000, fw_data)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Done flashing!")

	if *run {
		log.Println("Running binary from flash.")
		err = boot(port)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *echo {
		if *echo_baud != *baud {
			log.Printf("Switching from %d baud to %d baud.\n", *baud, *echo_baud)
			options.BaudRate = *echo_baud
			err := port.Close()
			if err != nil {
				log.Fatal(err)
			}
			port, err = serial.Open(options)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("Echoing from board...")
		fmt.Println("---------------------")
		carriage_return := false
		for true {
			buf := make([]byte, 1)
			n, err := port.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			if n > 0 {
				if buf[0] == '\r' {
					fmt.Printf("\n")
					carriage_return = true
					continue
				}
				if !(carriage_return && buf[0] == '\n') {
					fmt.Printf("%c", buf[0])
				}
				carriage_return = false
			}
		}
	}
}
