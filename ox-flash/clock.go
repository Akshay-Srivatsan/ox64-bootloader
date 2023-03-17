package main

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
)

const ClockMagic uint32 = 0x47464350

type ClockData struct {
	Magic  uint32
	Config ClockConfig
	Crc32  uint32
}

func NewClockData(config ClockConfig) ClockData {
	data := ClockData{
		Magic:  ClockMagic,
		Config: config,
		Crc32:  0,
	}
	var buf = bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, &config)
	data.Crc32 = crc32.ChecksumIEEE(buf.Bytes())
	return data
}

type ClockConfig struct {
	XtalType     uint8
	McuClk       uint8
	McuClkDiv    uint8
	McuBclkDiv   uint8
	McuPbclkDiv  uint8
	LpDiv        uint8
	DspClk       uint8
	DspClkDiv    uint8
	DspBclkDiv   uint8
	DspPbclk     uint8
	DspPbclkDiv  uint8
	EmiClk       uint8
	EmiClkDiv    uint8
	FlashClkType uint8
	FlashClkDiv  uint8
	WifipllPu    uint8
	AupllPu      uint8
	CpupllPu     uint8
	MipipllPu    uint8
	UhspllPu     uint8
}

var DefaultClockConfig = ClockConfig{
	XtalType:     4,
	McuClk:       4,
	McuClkDiv:    0,
	McuBclkDiv:   0,
	McuPbclkDiv:  3,
	LpDiv:        1,
	DspClk:       3,
	DspClkDiv:    0,
	DspBclkDiv:   1,
	DspPbclk:     2,
	DspPbclkDiv:  0,
	EmiClk:       2,
	EmiClkDiv:    1,
	FlashClkType: 1,
	FlashClkDiv:  0,
	WifipllPu:    1,
	AupllPu:      1,
	CpupllPu:     1,
	MipipllPu:    1,
	UhspllPu:     1,
}
