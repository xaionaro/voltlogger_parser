This code is supposed to parse data of voltlogger and [voltloggerA](https://devel.mephi.ru/dyokunev/voltloggerA).

For example:

    socat -u udp-recv:30319 - | ./voltlogger_parser -i - -n1 > ~/voltlogger.csv
    ^C
    qtiplot ~/voltlogger.csv
