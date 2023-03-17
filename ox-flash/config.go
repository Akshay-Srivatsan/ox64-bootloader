package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"hash/crc32"
	"log"
)

const HeaderMagic uint32 = 0x504e4642
const BootConfig uint32 = 0x654c0100

type BootHeader struct {
	MagicCode         uint32
	Revision          uint32
	Flash             FlashData
	Clock             ClockData
	BootConfig        uint32
	GroupImageOffset  uint32
	AesRegionLen      uint32
	ImgLenCnt         uint32
	Hash              [32]byte
	M0Config          CpuConfig
	D0Config          CpuConfig
	LpConfig          CpuConfig
	Boot2PtTable0     uint32
	Boot2PtTable1     uint32
	FlashCfgTableAddr uint32
	FlashCfgTableLen  uint32
	PatchRead         [4]PatchEntry
	PatchJump         [4]PatchEntry
	Reserved          [5]uint32
}

type BootInfo struct {
	Header          BootHeader
	Crc32           uint32
	FlashTableMagic uint32
	FlashTableCksum uint32
}

type CpuConfig struct {
	ConfigEnable       uint8
	HaltCpu            uint8
	CacheConfig        uint8
	Reserved           uint8
	CacheRangeH        uint32
	CacheRangeL        uint32
	ImageAddressOffset uint32
	BootEntry          uint32
	MspVal             uint32
}

type PatchEntry struct {
	Address uint32
	Value   uint32
}

func CreateBootHeader(jedec uint32, fw []byte, addr uint32) (BootHeader, error) {
	var header = BootHeader{}
	header.MagicCode = HeaderMagic
	header.Revision = 1

	header.Flash = NewFlashData(DefaultFlashConfig)
	header.Clock = NewClockData(DefaultClockConfig)

	header.Hash = sha256.Sum256(fw)

	header.BootConfig = BootConfig

	header.M0Config.ConfigEnable = 0
	header.M0Config.BootEntry = 0x58000000 + addr

	header.D0Config.ConfigEnable = 1
	header.D0Config.BootEntry = 0x58000000 + addr

	header.LpConfig.ConfigEnable = 0
	header.LpConfig.BootEntry = 0x58000000 + addr

	header.ImgLenCnt = uint32(len(fw))
	header.GroupImageOffset = 0x2000

	header.FlashCfgTableAddr = 352
	header.FlashCfgTableLen = 8

	return header, nil
}

func CreateBootInfo(jedec uint32, fw []byte, addr uint32) (BootInfo, error) {
	var info = BootInfo{}
	header, err := CreateBootHeader(jedec, fw, addr)
	if err != nil {
		return info, err
	}
	info.Header = header
	var hdr = bytes.NewBuffer([]byte{})
	binary.Write(hdr, binary.LittleEndian, &header)
	info.Crc32 = crc32.ChecksumIEEE(hdr.Bytes())
	log.Printf("Header: %d bytes\n", len(hdr.Bytes()))
	info.FlashTableMagic = 0x47544346
	info.FlashTableCksum = 0
	return info, nil
}
