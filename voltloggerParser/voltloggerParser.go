package voltloggerParser

import (
	"os"
	"io"
	"fmt"
	"strings"
	"encoding/binary"
)
const (
	WRITEBLOCK_SIZE	int64 = 512
)
/*
const (
	TIMESTAMP_WINDOW	int = 4096
)*/

type VoltloggerDumpRawHeader struct {
	Version			byte
	Magic			[11]byte
	Modificators		byte
	ChannelsNum		byte
	BlockWriteClockDelay	byte
	Reserved0		[1]byte
	DeviceName		[16]byte
	Reserved1		[480]byte
}

type VoltloggerDumpHeader struct {
	DeviceName		string
	NoClock			bool
	BlockWriteClockDelay	int64
	ChannelsNum		int;
}

func get16(dumpFile *os.File) (r uint16, err error) {
	err = binary.Read(dumpFile, binary.LittleEndian, &r)
	return r, err
}

func ParseVoltloggerDump(dumpPath string, noHeaders bool, channelsNum int, headerHandler func(VoltloggerDumpHeader, interface{})(error), rowHandler func(int64, []int32, VoltloggerDumpHeader, interface{})(error), arg interface{}) (err error) {
	var r VoltloggerDumpHeader

	// Openning the "dumpPath" as a file
	if (dumpPath == "-") {
		dumpPath = "/dev/stdin"
	}

	dumpFile, err := os.Open(dumpPath)
	if (err != nil) {
		return err
	}
	defer dumpFile.Close()

	if (noHeaders) {
		r.ChannelsNum = 1
		r.NoClock     = false
	} else {

		// Reading binary header to a VoltloggerDumpRawHeader
		var raw VoltloggerDumpRawHeader
		err = binary.Read(dumpFile, binary.LittleEndian, &raw)
		if (err != nil) {
			return fmt.Errorf("Cannot read dump: %v", err.Error())
		}

		// Checking if the data of a known type
		magicString := string(raw.Magic[:])
		if (magicString != "voltlogger\000") {
			return fmt.Errorf("The source is not a voltlogger dump (magic doesn't match: \"%v\" != \"voltlogger\")", magicString)
		}

		if (raw.Version != 0) {
			return fmt.Errorf("Unsupported dump version: %v", raw.Version)
		}
		if ((raw.Modificators & 0xfe) != 0) {
			return fmt.Errorf("Unsupported modificators bitmask: %o %o", raw.Modificators)
		}
		if (raw.ChannelsNum == 0) {
			return fmt.Errorf("Channels number is zero")
		}

		// Filling the VoltloggerDump struct
		r.DeviceName		= strings.Trim(string(raw.DeviceName[:]), "\000")
		r.NoClock		= (raw.Modificators & 0x01 != 0)
		r.BlockWriteClockDelay	= int64(raw.BlockWriteClockDelay)

		err = headerHandler(r, arg);
		if (err != nil) {
			return err;
		}

		r.ChannelsNum = int(raw.ChannelsNum)
	}

	if (r.ChannelsNum > 0) {
		r.ChannelsNum = channelsNum
	}

	// Parsing the Data

	var pos int64
	var timestampGlobal   int64
	var timestampLocalOld uint16
	timestampGlobal = -1

	for err = nil; err == nil; pos++ {

		if (r.NoClock) {
			timestampGlobal++
		} else {
			timestampLocal, err := get16(dumpFile)
			if (err != nil) {
				break
			}

			if (timestampGlobal < 0) {
				timestampGlobal = int64(timestampLocal)
			} else {
				var timestampLocalDiff int
				timestampLocalDiff = int(timestampLocal) - int(timestampLocalOld)
				/*if (timestampLocalDiff*timestampLocalDiff > TIMESTAMP_WINDOW*TIMESTAMP_WINDOW) {
					break
				}*/

				timestampLocalOld  = uint16(timestampLocal)
				if (timestampLocalDiff < 0) {
					timestampLocalDiff += (1 << 16)
				}

				timestampGlobal += int64(timestampLocalDiff)

			}
		}

		if (pos % WRITEBLOCK_SIZE == 0) {
			timestampGlobal += r.BlockWriteClockDelay
		}

		row := make([]int32, r.ChannelsNum)
		for i:=0; i < r.ChannelsNum; i++ {
			value, err := get16(dumpFile)
			if (err != nil) {
				break
			}
			row[i] = int32(value)
		}

		err = rowHandler(timestampGlobal, row, r, arg);
		if (err != nil) {
			return fmt.Errorf("Got an error from rowHandler: %v", err);
		}
	}

	if (err == io.EOF) {
		err = nil
	}

	return err
}

