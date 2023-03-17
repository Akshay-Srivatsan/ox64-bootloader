package main

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
)

const FlashMagic uint32 = 0x47464346

type FlashData struct {
	Magic  uint32
	Config FlashConfig
	Crc32  uint32
}

func NewFlashData(config FlashConfig) FlashData {
	data := FlashData{
		Magic:  FlashMagic,
		Config: config,
		Crc32:  0,
	}
	var buf = bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, &config)
	data.Crc32 = crc32.ChecksumIEEE(buf.Bytes())
	return data
}

type FlashConfig struct {
	IoMode               uint8
	ContReadSupport      uint8
	SfctrlClkDelay       uint8
	SfctrlClkInvert      uint8
	ResetEnCmd           uint8
	ResetCmd             uint8
	ExitContreadCmd      uint8
	ExitContreadCmdSize  uint8
	JedecidCmd           uint8
	JedecidCmdDmyClk     uint8
	Enter32BitsAddrCmd   uint8
	Exit32BitsAddrClk    uint8
	SectorSize           uint8
	MfgId                uint8
	PageSize             uint16
	ChipEraseCmd         uint8
	SectorEraseCmd       uint8
	Blk32KEraseCmd       uint8
	Blk64KEraseCmd       uint8
	WriteEnableCmd       uint8
	PageProgCmd          uint8
	QpageProgCmd         uint8
	QualPageProgAddrMode uint8
	FastReadCmd          uint8
	FastReadDmyClk       uint8
	QpiFastReadCmd       uint8
	QpiFastReadDmyClk    uint8
	FastReadDoCmd        uint8
	FastReadDoDmyClk     uint8
	FastReadDioCmd       uint8
	FastReadDioDmyClk    uint8
	FastReadQoCmd        uint8
	FastReadQoDmyClk     uint8
	FastReadQioCmd       uint8
	FastReadQioDmyClk    uint8
	QpiFastReadQioCmd    uint8
	QpiFastReadQioDmyClk uint8
	QpiPageProgCmd       uint8
	WriteVregEnableCmd   uint8
	WelRegIndex          uint8
	QeRegIndex           uint8
	BusyRegIndex         uint8
	WelBitPos            uint8
	QeBitPos             uint8
	BusyBitPos           uint8
	WelRegWriteLen       uint8
	WelRegReadLen        uint8
	QeRegWriteLen        uint8
	QeRegReadLen         uint8
	ReleasePowerDown     uint8
	BusyRegReadLen       uint8
	RegReadCmd0          uint8
	RegReadCmd1          uint8
	ShouldBeZero0        uint16
	RegWriteCmd0         uint8
	RegWriteCmd1         uint8
	ShouldBeZero1        uint16
	EnterQpiCmd          uint8
	ExitQpiCmd           uint8
	ContReadCode         uint8
	ContReadExitCode     uint8
	BurstWrapCmd         uint8
	BurstWrapDmyClk      uint8
	BurstWrapDataMode    uint8
	BurstWrapCode        uint8
	DeBurstWrapCmd       uint8
	DeBurstWrapCmdDmyClk uint8
	DeBurstWrapCodeMode  uint8
	DeBurstWrapCode      uint8
	SectorEraseTime      uint16
	Blk32KEraseTime      uint16
	Blk64KEraseTime      uint16
	PageProgTime         uint16
	ChipEraseTime        uint16
	PowerDownDelay       uint8
	QeData               uint8
}

var DefaultFlashConfig = FlashConfig{
	IoMode:               16,
	ContReadSupport:      0,
	SfctrlClkDelay:       1,
	SfctrlClkInvert:      1,
	ResetEnCmd:           102,
	ResetCmd:             153,
	ExitContreadCmd:      255,
	ExitContreadCmdSize:  3,
	JedecidCmd:           159,
	JedecidCmdDmyClk:     0,
	Enter32BitsAddrCmd:   183,
	Exit32BitsAddrClk:    233,
	SectorSize:           4,
	MfgId:                255,
	PageSize:             256,
	ChipEraseCmd:         199,
	SectorEraseCmd:       32,
	Blk32KEraseCmd:       82,
	Blk64KEraseCmd:       216,
	WriteEnableCmd:       6,
	PageProgCmd:          2,
	QpageProgCmd:         50,
	QualPageProgAddrMode: 0,
	FastReadCmd:          11,
	FastReadDmyClk:       1,
	QpiFastReadCmd:       11,
	QpiFastReadDmyClk:    1,
	FastReadDoCmd:        59,
	FastReadDoDmyClk:     1,
	FastReadDioCmd:       187,
	FastReadDioDmyClk:    0,
	FastReadQoCmd:        107,
	FastReadQoDmyClk:     1,
	FastReadQioCmd:       235,
	FastReadQioDmyClk:    2,
	QpiFastReadQioCmd:    235,
	QpiFastReadQioDmyClk: 2,
	QpiPageProgCmd:       2,
	WriteVregEnableCmd:   80,
	WelRegIndex:          0,
	QeRegIndex:           1,
	BusyRegIndex:         0,
	WelBitPos:            1,
	QeBitPos:             1,
	BusyBitPos:           0,
	WelRegWriteLen:       2,
	WelRegReadLen:        1,
	QeRegWriteLen:        2,
	QeRegReadLen:         1,
	ReleasePowerDown:     171,
	BusyRegReadLen:       1,
	RegReadCmd0:          5,
	RegReadCmd1:          53,
	ShouldBeZero0:        0,
	RegWriteCmd0:         1,
	RegWriteCmd1:         1,
	ShouldBeZero1:        0,
	EnterQpiCmd:          56,
	ExitQpiCmd:           255,
	ContReadCode:         32,
	ContReadExitCode:     240,
	BurstWrapCmd:         119,
	BurstWrapDmyClk:      3,
	BurstWrapDataMode:    2,
	BurstWrapCode:        64,
	DeBurstWrapCmd:       119,
	DeBurstWrapCmdDmyClk: 3,
	DeBurstWrapCodeMode:  2,
	DeBurstWrapCode:      240,
	SectorEraseTime:      300,
	Blk32KEraseTime:      1200,
	Blk64KEraseTime:      1200,
	PageProgTime:         50,
	ChipEraseTime:        33000,
	PowerDownDelay:       20,
	QeData:               0,
}
