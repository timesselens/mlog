package main

import (
	"control"
	"flag"
	"fmt"
	"github.com/robfig/config"
    "mlog/core"
    "mlog/webserver"
	"os"
	"strings"
)

var HTTP_LISTEN string
var WEBSOCKET_LISTEN string
var SYSLOG_LISTEN string
var PIPE_CONTROL_DIR string

func init() {
	flag.StringVar(&HTTP_LISTEN, "http", "127.0.0.1:1980", "webserver listen host and port (tcp)")
	flag.StringVar(&WEBSOCKET_LISTEN, "ws", "127.0.0.1:12345", "websocket listen host and port (tcp)")
	flag.StringVar(&SYSLOG_LISTEN, "syslog", "0.0.0.0:9999", "syslog listen host and port (udp)")
	flag.StringVar(&PIPE_CONTROL_DIR, "controldir", "var", "pipe control dir")
	flag.BoolVar(&core.Config.Verbose, "verbose", false, "be verbose")

	tmpDebug := ""
	flag.StringVar(&tmpDebug, "debug", "", "debug a particular part of the code")
	for _, s := range strings.Split(tmpDebug, " ") {
		switch s {
		case "core":
			core.Config.Debug.Core = true
        case "http":
            core.Config.Debug.HTTP = true
        case "process":
            core.Config.Debug.Process = true
		}
	}

	flag.Usage = func() {
		var help = `
mlog v.0.0.1

 -http       127.0.0.1:1980      webserver listen host and port (tcp)
 -ws         127.0.0.1:12345     websocket listen host and port (tcp)
 -syslog     0.0.0.0:9999        syslog listen host and port (udp)
 -controldir var/                directory containing named pipes
 -verbose    false               be verbose
 -debug      ""                  debug a particular part of the code (comma sep)
 -V                              show version information

examples:

 ./mlog -controldir /var/run -debug process,http
        `
		fmt.Println(help)
		os.Exit(0)
	}
}

func main() {
	ver := flag.Bool("V", false, "print version information")
	flag.Parse()

	if *ver {
		fmt.Println("mlog v0.0.1")
		os.Exit(0)
	}
	fmt.Println("starting mlog")
	c, _ := config.ReadDefault("etc/mlog.cfg")
	fmt.Println("read config", c)

    go control.UnixDomainServer(PIPE_CONTROL_DIR)

    webserver.HTTPServer(HTTP_LISTEN)
}
