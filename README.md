This code is supposed to parse data of voltlogger and [voltloggerA](https://devel.mephi.ru/dyokunev/voltloggerA).

The main repository: git clone [https://devel.mephi.ru/dyokunev/voltlogger_parser](https://devel.mephi.ru/dyokunev/voltlogger_parser)

Compiling:

    export GOPATH=$HOME:/gocode
    go get devel.mephi.ru/dyokunev/voltlogger_parser
    cd "$GOPATH/src/devel.mephi.ru/dyokunev/voltlogger_parser"
    make
    ./voltlogger_parser --help

Running example:

    socat -u udp-recv:30319 - | ./voltlogger_parser/voltlogger_parser -i - -n > ~/voltlogger.csv
    ^C
    qtiplot ~/voltlogger.csv

Another example:

    socat -u udp-recv:30319 - | ./voltlogger_parser/voltlogger_parser -b -i - -n -t > ~/voltage.binlog &
    ./voltlogger_oscilloscope/voltlogger_oscilloscope -i ~/voltage.binlog -t

[https://devel.mephi.ru/dyokunev/voltlogger_oscilloscope](https://devel.mephi.ru/dyokunev/voltlogger_oscilloscope)

