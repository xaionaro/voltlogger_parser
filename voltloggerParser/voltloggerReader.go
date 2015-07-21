package voltloggerReader

import (
	"os"
	"fmt"
	"encoding/binary"
)

type VoltloggerDumpRaw struct {
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
	Data		map[int]map[int]int
}

func ParseVoltloggerDump(dumpPath string) (r VoltloggerDump, err error) {
	// Openning the "dumpPath" as a file
	dumpFile, err := os.Open(dumpPath)
	if (err != nil) {
		return r, err
	}
	defer dumpFile.Close()

	// Reading binary data to a VoltloggerDumpRaw
	var raw VoltloggerDumpRaw
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



	// Filling the VoltloggerDump struct
	r.DeviceName = string(raw.DeviceName[:])

	return r, nil
}
