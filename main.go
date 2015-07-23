package main

import (
	"os"
	"fmt"
	"time"
	"code.google.com/p/getopt"
	"devel.mephi.ru/dyokunev/voltlogger_parser/voltloggerParser"
)

type printRow_arg struct {
	outputPath	string
}

func handleHeader(h voltloggerParser.VoltloggerDumpHeader, arg_iface interface{}) (error) {
	arg := arg_iface.(*printRow_arg)

	if (arg.outputPath == "") {
		now              := time.Now()
		year, month, day := now.Date()
		hour, min,   sec := now.Clock()
		arg.outputPath    = fmt.Sprintf("%v_%v-%02v-%02v_%02v:%02v:%02v.csv", h.DeviceName, year, int(month), day, hour, min, sec)
	}

	return nil
}

func printRow(ts int64, row []int, h voltloggerParser.VoltloggerDumpHeader, arg_iface interface{}) (error) {
	fmt.Printf("%v", ts)
	rowLen := len(row)
	for i:=0; i < rowLen; i++ {
		fmt.Printf("\t%v", row[i])
	}
	fmt.Printf("\n")
	return nil
}

func main() {
	var dumpPath		string
	var noHeadersStr	string
	var noHeaders		bool
	var printRow_arg	printRow_arg

	getopt.StringVar(&dumpPath,			'i',	"dump-path")
	getopt.StringVar(&printRow_arg.outputPath,	'o',	"output-path").SetOptional()
	getopt.StringVar(&noHeadersStr,			'n',	"no-headers").SetOptional()

	getopt.Parse()
	if (getopt.NArgs() > 0 || dumpPath == "") {
		getopt.Usage()
		os.Exit(-2)
	}

	switch (noHeadersStr) {
		case "1", "yes", "true":
			noHeaders = true
			break
	}

	err := voltloggerParser.ParseVoltloggerDump(dumpPath, noHeaders, handleHeader, printRow, &printRow_arg)
	if (err != nil) {
		fmt.Printf("Cannot parse the dump: %v\n", err.Error())
		os.Exit(-1)
	}

	fmt.Printf("%v %v\n", dumpPath, printRow_arg)
}
