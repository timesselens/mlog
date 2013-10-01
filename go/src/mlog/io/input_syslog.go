package io

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type SyslogListener struct {
    ID
    Hostport string
    LastError string
    *stdchannels
	addr *net.UDPAddr
    sock *net.UDPConn
}

// implement Stringer interface
func (sl SyslogListener) String() (string) {
    return "syslog-"+string(sl.ID) +"://" + sl.Hostport + "#" + strings.Replace(sl.LastError," ","_",-1)
}

func NewSyslogListener(hostport string) (sl *SyslogListener, err error) {
	sl = new(SyslogListener)
    sl.Hostport = hostport
    err = sl.Init()
	return sl,err
}

func (sl *SyslogListener) Init() (err error) {
	addr, _ := net.ResolveUDPAddr("udp", sl.Hostport)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to convert hostport into address: %s", sl.Hostport))
	}
    sl.ID = ID("syslog@"+sl.Hostport)
    sl.stdchannels = &stdchannels{}
    sl.in = make(chan string)
    sl.out = make(chan string)
    sl.addr = addr
	go sl.Handler()
    return nil
}

func (sl *SyslogListener) Handler() {
    sock, err := net.ListenUDP("udp", sl.addr)
    sl.sock = sock
	defer sock.Close()

	if err != nil {
        sl.LastError = err.Error()
        log.Printf("error setting up new SyslogListener: %s", err)
	}

	buffer := bufio.NewReader(sock)
	for {
		line, err := buffer.ReadString('\n')
		if err != nil {
			break
		}
		sl.out <- strings.TrimRight(line, "\n")
	}
    log.Printf("stopping sysloglistener")
}

