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
var keep_open = flag.Bool("keep-open", false, "keep the serial port open after flashing")
var payload = flag.Bool("payload", false, "flash a payload (not firmware)")
var addr = flag.Uint("addr", 0, "the address to flash to (0x0 for firmware, 0x10000 for a payload)")

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

	if *keep_open && *payload {
		log.Fatal("the --keep-open and --payload")
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

	fmt.Printf("Connected to %s, baud: %d\n", options.PortName, options.BaudRate)
	if *keep_open {
		fmt.Println("Keeping connection open after flashing.")
	}

	jedec, err := handshake(port, options.BaudRate)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Have %s file: %s\n", filetype, fw)
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

		fmt.Println("Writing bootrom-compatible header at 0x0000")
		err = flash(port, 0, bh.Bytes())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Writing firmware at 0x%x\n", *addr)
		err = flash(port, uint32(*addr)+0x2000, fw_data)
		if err != nil {
			log.Fatal(err)
		}

		err = boot(port)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Writing payload at 0x%x\n", *addr)
		err = flash(port, uint32(*addr)+0x2000, fw_data)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Done flashing!")

	if *keep_open {
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
