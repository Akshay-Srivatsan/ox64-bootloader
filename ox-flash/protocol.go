package main

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"reflect"
	"time"
)

const (
	GetBootInfo         byte = 0x10
	ClockSetPara             = 0x22
	FlashErase               = 0x30
	FlashWrite               = 0x31
	FlashRead                = 0x32
	FlashBoot                = 0x33
	FlashReadJedecId         = 0x36
	FlashReadStatusReg       = 0x37
	FlashWriteStatusReg      = 0x38
	FlashWriteCheck          = 0x3A
	FlashSetPara             = 0x3B
	FlashChipErase           = 0x3C
	FlashReadSha             = 0x3D
	FlashXipReadSha          = 0x3E
	EfuseReadMacAddress      = 0x42
	MemWrite                 = 0x50
	XipReadStart             = 0x60
	XipReadFinish            = 0x61
	LogRead                  = 0x71
)

func sendTwoWordMessage(port io.ReadWriter, command byte, a uint32, b uint32) error {
	buf := make([]byte, 0, 8)
	buf = binary.LittleEndian.AppendUint32(buf, a)
	buf = binary.LittleEndian.AppendUint32(buf, b)
	return sendMessage(port, command, buf)
}

func sendTwoWordDataMessage(port io.ReadWriter, command byte, a uint32, b uint32, data []byte) error {
	buf := make([]byte, 0, 8)
	buf = binary.LittleEndian.AppendUint32(buf, a)
	buf = binary.LittleEndian.AppendUint32(buf, b)
	return sendMessage(port, command, append(buf, data...))
}

func sendWordMessage(port io.ReadWriter, command byte, a uint32) error {
	buf := make([]byte, 0, 4)
	buf = binary.LittleEndian.AppendUint32(buf, a)
	return sendMessage(port, command, buf)
}

func sendWordDataMessage(port io.ReadWriter, command byte, a uint32, buf []byte) error {
	var x [4]byte
	binary.LittleEndian.PutUint32(x[:], a)

	return sendMessage(port, command, append(x[:], buf...))
}

func sendEmptyMessage(port io.ReadWriter, command byte) error {
	buf := make([]byte, 0, 0)
	return sendMessage(port, command, buf)
}

func sendMessage(port io.ReadWriter, command byte, data []byte) error {
	log.Printf("Sending command 0x%x with %d bytes of data", command, len(data))
	buf := make([]byte, len(data)+4)
	var bytes [2]byte
	binary.LittleEndian.PutUint16(bytes[:], uint16(len(data)))

	var cksum uint8 = 0
	buf[0] = command
	buf[1] = 0
	buf[2] = bytes[0]
	buf[3] = bytes[1]

	cksum += bytes[0]
	cksum += bytes[1]
	for i, v := range data {
		cksum += v
		buf[i+4] = v
	}

	buf[1] = cksum

	// for i, v := range buf {
	// 	log.Printf("%d = 0x%x\n", i, v)
	// }

	attempt := 0
	for attempt < 10 {
		attempt++
		// port.Write(buf)
		err := writeExact(port, buf)
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		var response [2]byte
		err = readExact(port, response[:])
		if err != nil {
			return err
		}
		log.Printf("Response: [%x, %x]", response[0], response[1])
		if "OK" == string(response[:]) {
			log.Printf("Saw OK response")
			return nil
		} else if "FL" == string(response[:]) {
			err := readExact(port, response[:])
			if err != nil {
				return err
			}
			switch binary.LittleEndian.Uint16(response[:]) {
			case 0x0000:
				return nil
			/*flash*/
			case 0x0001:
				return errors.New("FLASH_INIT_ERROR")
			case 0x0002:
				return errors.New("FLASH_ERASE_PARA_ERROR")
			case 0x0003:
				return errors.New("FLASH_ERASE_ERROR")
			case 0x0004:
				return errors.New("FLASH_WRITE_PARA_ERROR")
			case 0x0005:
				return errors.New("FLASH_WRITE_ADDR_ERROR")
			case 0x0006:
				return errors.New("FLASH_WRITE_ERROR")
			case 0x0007:
				return errors.New("FLASH_BOOT_PARA_ERROR")
			case 0x0008:
				return errors.New("FLASH_SET_PARA_ERROR")
			case 0x0009:
				return errors.New("FLASH_READ_STATUS_REG_ERROR")
			case 0x000A:
				return errors.New("FLASH_WRITE_STATUS_REG_ERROR")
			/*cmd*/
			case 0x0101:
				return errors.New("CMD_ID_ERROR ")
			case 0x0102:
				return errors.New("CMD_LEN_ERROR")
			case 0x0103:
				return errors.New("CMD_CRC_ERROR")
			case 0x0104:
				return errors.New("CMD_SEQ_ERROR")
			/*image*/
			case 0x0201:
				return errors.New("IMG_BOOTHEADER_LEN_ERROR")
			case 0x0202:
				return errors.New("IMG_BOOTHEADER_NOT_LOAD_ERROR")
			case 0x0203:
				return errors.New("IMG_BOOTHEADER_MAGIC_ERROR")
			case 0x0204:
				return errors.New("IMG_BOOTHEADER_CRC_ERROR")
			case 0x0205:
				return errors.New("IMG_BOOTHEADER_ENCRYPT_NOTFIT")
			case 0x0206:
				return errors.New("IMG_BOOTHEADER_SIGN_NOTFIT")
			case 0x0207:
				return errors.New("IMG_SEGMENT_CNT_ERROR")
			case 0x0208:
				return errors.New("IMG_AES_IV_LEN_ERROR")
			case 0x0209:
				return errors.New("IMG_AES_IV_CRC_ERROR")
			case 0x020a:
				return errors.New("IMG_PK_LEN_ERROR")
			case 0x020b:
				return errors.New("IMG_PK_CRC_ERROR")
			case 0x020c:
				return errors.New("IMG_PK_HASH_ERROR")
			case 0x020d:
				return errors.New("IMG_SIGNATURE_LEN_ERROR")
			case 0x020e:
				return errors.New("IMG_SIGNATURE_CRC_ERROR")
			case 0x020f:
				return errors.New("IMG_SECTIONHEADER_LEN_ERROR")
			case 0x0210:
				return errors.New("IMG_SECTIONHEADER_CRC_ERROR")
			case 0x0211:
				return errors.New("IMG_SECTIONHEADER_DST_ERROR")
			case 0x0212:
				return errors.New("IMG_SECTIONDATA_LEN_ERROR")
			case 0x0213:
				return errors.New("IMG_SECTIONDATA_DEC_ERROR")
			case 0x0214:
				return errors.New("IMG_SECTIONDATA_TLEN_ERROR")
			case 0x0215:
				return errors.New("IMG_SECTIONDATA_CRC_ERROR")
			case 0x0216:
				return errors.New("IMG_HALFBAKED_ERROR")
			case 0x0217:
				return errors.New("IMG_HASH_ERROR")
			case 0x0218:
				return errors.New("IMG_SIGN_PARSE_ERROR")
			case 0x0219:
				return errors.New("IMG_SIGN_ERROR")
			case 0x021a:
				return errors.New("IMG_DEC_ERROR")
			case 0x021b:
				return errors.New("IMG_ALL_INVALID_ERROR")
			/*IF*/
			case 0x0301:
				return errors.New("IF_RATE_LEN_ERROR")
			case 0x0302:
				return errors.New("IF_RATE_PARA_ERROR")
			case 0x0303:
				return errors.New("IF_PASSWORDERROR")
			case 0x0304:
				return errors.New("IF_PASSWORDCLOSE")
			/*MISC*/
			case 0xfffc:
				return errors.New("PLL_ERROR")
			case 0xfffd:
				return errors.New("INVASION_ERROR")
			case 0xfffe:
				return errors.New("POLLING")
			case 0xffff:
				return errors.New("FAIL")
			}
		}
		log.Printf("invalid response %x; retrying", response)
		var discard [1024]byte
		port.Read(discard[:])
	}
	return errors.New("invalid response from Ox64")
}

func writeExact(port io.Writer, buf []byte) error {
	log.Printf("Writing %d bytes", len(buf))
	n := 0
	for n < len(buf) {
		x, err := port.Write(buf[n:])
		if err != nil {
			return err
		}
		n += x
	}
	return nil
}

func readExact(port io.Reader, buf []byte) error {
	log.Printf("Reading %d bytes", len(buf))
	n := 0
	for n < len(buf) {
		x, err := port.Read(buf[n:])
		if err != nil {
			return err
		}
		n += x
	}
	return nil
}

func readResponseData(port io.Reader) ([]byte, error) {
	var length [2]byte
	err := readExact(port, length[:])
	if err != nil {
		return nil, err
	}
	buf := make([]byte, binary.LittleEndian.Uint16(length[:]))
	err = readExact(port, buf[:])
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func handshake(port io.ReadWriter, baud uint) (uint32, error) {
	// transmit the baud autodetection sequence (at least 5ms)
	n := baud*5/10000 + 1
	log.Printf("Running sync sequence (%d bytes)\n", n)
	sync := make([]byte, n)
	for i := range sync {
		sync[i] = 0x55
	}
	writeExact(port, sync)

	time.Sleep(30 * time.Millisecond)

	// reset boot status register or something (not super clear what this does)
	err := sendTwoWordMessage(port, MemWrite, 0x2000F038, 0x18000000)
	if err != nil {
		return 0, err
	}
	log.Println("Done with handshake")

	log.Println("Getting JEDEC ID")
	err = sendWordMessage(port, FlashSetPara, 0x14180)
	if err != nil {
		return 0, err
	}
	err = sendEmptyMessage(port, FlashReadJedecId)
	if err != nil {
		return 0, err
	}
	data, err := readResponseData(port)
	if err != nil {
		return 0, err
	}
	jedec := binary.LittleEndian.Uint32(data)
	log.Printf("Got JEDEC ID 0x%x\n", jedec)
	return jedec, nil
}

func flash(port io.ReadWriter, addr uint32, data []byte) error {
	log.Println("Setting flash parameters for write")
	err := sendWordDataMessage(port, FlashSetPara, 0x14180, flash_para)
	if err != nil {
		return err
	}

	n := len(data)
	end := int(addr) + n
	log.Printf("Erasing from 0x%x to 0x%x", addr, end-1)
	err = sendTwoWordMessage(port, FlashErase, addr, uint32(end-1))
	if err != nil {
		return err
	}

	log.Printf("Writing data to flash (%d bytes)", n)
	packetSize := 4096 - 8
	packets := (n / packetSize) + 1
	if packets > 1 {
		log.Printf("Using %d packets", packets)
	}
	for i := 0; i < packets; i++ {
		start := i * packetSize
		end := start + packetSize
		if end > n {
			end = n
		}
		packetData := data[start:end]
		log.Printf("Writing packet %d (%d bytes) to 0x%x\n", i, end-start, start)
		err = sendWordDataMessage(port, FlashWrite, addr+uint32(start), packetData)
		if err != nil {
			return err
		}
	}

	log.Println("Running write check")
	err = sendEmptyMessage(port, FlashWriteCheck)
	if err != nil {
		return err
	}

	log.Println("Switching into XIP mode")
	err = sendEmptyMessage(port, XipReadStart)
	if err != nil {
		return err
	}

	log.Println("Reading SHA256")
	err = sendTwoWordMessage(port, FlashXipReadSha, addr, uint32(n))
	if err != nil {
		return err
	}
	resp, err := readResponseData(port)
	if err != nil {
		return err
	}

	sum := sha256.Sum256(data)
	log.Printf("expect SHA256 is %x", sum)
	log.Printf("device SHA256 is %x", resp)

	if !reflect.DeepEqual(sum[:], resp) {
		return errors.New("SHA256 did not match!")
	}
	return nil
}

func readFlash(port io.ReadWriter, addr uint32, len uint32) ([]byte, error) {
	log.Println("Setting flash parameters for read")
	err := sendWordDataMessage(port, FlashSetPara, 0x14180, flash_para)
	if err != nil {
		return nil, err
	}

	err = sendTwoWordMessage(port, FlashRead, addr, len)
	if err != nil {
		return nil, err
	}
	return readResponseData(port)
}

func boot(port io.ReadWriter) error {
	log.Println("Booting")
	err := sendWordMessage(port, FlashBoot, 0x58000000)
	return err
}
