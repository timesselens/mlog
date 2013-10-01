package io

import (
	"fmt"
	"os"
)

type Debug struct {
    ID
	Prefix string
    *stdchannels
}

func NewDebug(prefix string) (debug *Debug, err error) {
	debug = new(Debug)
    debug.Prefix = prefix
    err = debug.Init()
	return debug,err
}

func (debug *Debug) Init() (err error) {
    if len(debug.ID) == 0 {
        debug.ID = ID("debug://"+debug.Prefix+"@")
    } else {
        debug.ID = ID("debug://"+string(debug.ID)+"@"+debug.Prefix)
    }
    debug.stdchannels = &stdchannels{}
    debug.in = make(chan string)
    debug.exit = make(chan struct{})
	go debug.Handler()
    return nil
}

func (debug *Debug) Handler() {
LOOP:
	for {
		select {
		case line, ok := <-debug.in:
            fmt.Println("<<< ", line)
			if !ok {
				fmt.Println("channel closed, exitting")
				break LOOP
			}
			fmt.Fprintf(os.Stderr, "%s %s\n", debug.Prefix, line)
        case <-debug.exit:
            break LOOP;
		}
	}
	fmt.Println("debug handler finished")
}
