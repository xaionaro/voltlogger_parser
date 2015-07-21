package main

import (
	"os"
	"fmt"
	"time"
	"code.google.com/p/getopt"
	"devel.mephi.ru/dyokunev/voltlogger_parser/voltloggerParser"
)

func main() {
	var dumpPath	string
	var outputPath	string

	getopt.StringVar(&dumpPath,	'i',	"dump-path")
	getopt.StringVar(&outputPath,	'o',	"output-path").SetOptional()

	getopt.Parse()
	if (getopt.NArgs() > 0 || dumpPath == "") {
		getopt.Usage()
		os.Exit(-2)
	}

	dump, err := voltloggerReader.ParseVoltloggerDump(dumpPath)
	if (err != nil) {
		fmt.Printf("Cannot parse the dump: %v\n", err.Error())
		os.Exit(-1)
	}

	if (outputPath == "") {
		now              := time.Now()
		year, month, day := now.Date()
		hour, min,   sec := now.Clock()
		outputPath        = fmt.Sprintf("%v_%v-%02v-%02v_%02v:%02v:%02v.csv", dump.DeviceName, year, int(month), day, hour, min, sec)
	}

	fmt.Printf("%v %v %v\n", dumpPath, outputPath, dump)
}
