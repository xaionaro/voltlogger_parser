package voltloggerReader

import (
	"os"
	"io"
	"fmt"
	"encoding/binary"
)
/*
const (
	TIMESTAMP_WINDOW	int = 4096
)*/

type VoltloggerDumpRawHeader struct {
	Version		byte
	Magic		[11]byte
	Modificators	byte
	ChannelsNum	byte
	Reserved0	[2]byte
	DeviceName	[16]byte
	Reserved1	[480]byte
}

type VoltloggerDump struct {
	DeviceName	string
	Data		map[int64][]int
}

func get16(dumpFile *os.File) (r uint16, err error) {
	err = binary.Read(dumpFile, binary.LittleEndian, &r)
	return r, err
}

func ParseVoltloggerDump(dumpPath string) (r VoltloggerDump, err error) {
	// Openning the "dumpPath" as a file
	dumpFile, err := os.Open(dumpPath)
	if (err != nil) {
		return r, err
	}
	defer dumpFile.Close()

	// Reading binary header to a VoltloggerDumpRawHeader
	var raw VoltloggerDumpRawHeader
	err = binary.Read(dumpFile, binary.LittleEndian, &raw)
	if (err != nil) {
		return r, fmt.Errorf("Cannot read dump: %v", err.Error())
	}

	// Checking if the data of a known type
	magicString := string(raw.Magic[:])
	if (magicString != "voltlogger\000") {
		return r, fmt.Errorf("The source is not a voltlogger dump (magic doesn't match: \"%v\" != \"voltlogger\")", magicString)
	}

	if (raw.Version != 0) {
		return r, fmt.Errorf("Unsupported dump version: %v", raw.Version)
	}
	if (raw.Modificators != 0) {
		return r, fmt.Errorf("Unsupported modificators bitmask: %o", raw.Modificators)
	}
	if (raw.ChannelsNum == 0) {
		return r, fmt.Errorf("Channels number is zero")
	}

	// Parsing the Data

	channelsNum := int(raw.ChannelsNum)
	r.Data       = make(map[int64][]int)
	var timestampGlobal int64
	timestampGlobal = -1

	for err = nil; err == nil; {
		timestampLocal, err := get16(dumpFile)
		if (err != nil) {
			break
		}

		if (timestampGlobal < 0) {
			timestampGlobal = int64(timestampLocal)
		} else {
			var timestampLocalDiff int
			timestampLocalOld := int16(timestampGlobal)
			timestampLocalDiff = int(timestampLocal) - int(timestampLocalOld)
			/*if (timestampLocalDiff*timestampLocalDiff > TIMESTAMP_WINDOW*TIMESTAMP_WINDOW) {
				break
			}*/

			if (timestampLocalDiff <= 0) {
				timestampLocalDiff += (1 << 16)
			}

			timestampGlobal += int64(timestampLocalDiff);

		}

		values := make([]int, raw.ChannelsNum)
		for i:=0; i < channelsNum; i++ {
			value, err := get16(dumpFile)
			if (err != nil) {
				break
			}
			values[i] = int(value)
		}

		r.Data[timestampGlobal] = values;

	}

	if (err == io.EOF) {
		err = nil
	}

	// Filling the VoltloggerDump struct
	r.DeviceName = string(raw.DeviceName[:])

	return r, err
}
