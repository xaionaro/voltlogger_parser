package main

import (
	"os"
	"fmt"
	"time"
	"encoding/binary"
	"code.google.com/p/getopt"
	"devel.mephi.ru/dyokunev/voltlogger_parser/voltloggerParser"
)

type printRow_arg struct {
	outputPath	string
	outputFile	*os.File
	binaryOutput	bool
	insertParseTime	bool
}

func handleHeader(h voltloggerParser.VoltloggerDumpHeader, arg_iface interface{}) (err error) {
	arg := arg_iface.(*printRow_arg)
	if (arg.binaryOutput) {
		err = binary.Write(arg.outputFile, binary.LittleEndian, h)
		if (err != nil) {
			return err
		}
	}

	return nil
}

func printRow(ts int64, row []int32, h voltloggerParser.VoltloggerDumpHeader, arg_iface interface{}) (err error) {
	arg := arg_iface.(*printRow_arg)
	var parseTime int64

	if (arg.insertParseTime) {
		parseTime = time.Now().UnixNano()
	}

	if (arg.binaryOutput) {
		if (arg.insertParseTime) {
			err = binary.Write(arg.outputFile, binary.LittleEndian, parseTime)
			if (err != nil) {
				return err
			}
		}
		err = binary.Write(arg.outputFile, binary.LittleEndian, ts)
		if (err != nil) {
			return err
		}
		err = binary.Write(arg.outputFile, binary.LittleEndian, row)
		if (err != nil) {
			return err
		}
	} else {
		if (arg.insertParseTime) {
			fmt.Printf("%v\t", parseTime)
		}
		fmt.Printf("%v", ts)
		rowLen := len(row)
		for i:=0; i < rowLen; i++ {
			fmt.Printf("\t%v", row[i])
		}
		fmt.Printf("\n")
	}
	return nil
}

func main() {
	var err			error
	var dumpPath		string
	var noHeaders		bool
	var printRow_arg	printRow_arg
	var channelsNum		int

	getopt.StringVar(&dumpPath,			'i',	"dump-path"		)
	getopt.StringVar(&printRow_arg.outputPath,	'o',	"output-path"		).SetOptional()
	getopt.BoolVar  (&noHeaders,			'n',	"no-headers"		).SetOptional()
	getopt.IntVar   (&channelsNum,			'c',	"force-channels-num"	).SetOptional()
	getopt.BoolVar  (&printRow_arg.binaryOutput,	'b',	"binary-output"		).SetOptional()
	getopt.BoolVar  (&printRow_arg.insertParseTime,	't',	"insert-parse-time"	).SetOptional()

	getopt.Parse()
	if (getopt.NArgs() > 0 || dumpPath == "") {
		getopt.Usage()
		os.Exit(-2)
	}
	switch (printRow_arg.outputPath) {
/*		case "":
			now              := time.Now()
			year, month, day := now.Date()
			hour, min,   sec := now.Clock()
			printRow_arg.outputPath = fmt.Sprintf("%v_%v-%02v-%02v_%02v:%02v:%02v.csv", h.DeviceName, year, int(month), day, hour, min, sec)
			break*/
		case "", "-":
			printRow_arg.outputPath = "/dev/stdout"
			printRow_arg.outputFile = os.Stdout
			break
		default:
			printRow_arg.outputFile,err = os.Open(printRow_arg.outputPath)
			panic(fmt.Errorf("Not supported yet"))
	}

	if (err != nil) {
		fmt.Printf("Cannot open output file: %v", err.Error())
		os.Exit(-1)
	}
	//if (printRow_arg.binaryOutput) {
	//	err = binary.Write(printRow_arg.outputFile, binary.LittleEndian, printRow_arg)
	//}

	err = voltloggerParser.ParseVoltloggerDump(dumpPath, noHeaders, channelsNum, handleHeader, printRow, &printRow_arg)
	if (err != nil) {
		fmt.Printf("Cannot parse the dump: %v\n", err.Error())
		os.Exit(-1)
	}
	printRow_arg.outputFile.Close()

	fmt.Printf("%v %v\n", dumpPath, printRow_arg)
}
