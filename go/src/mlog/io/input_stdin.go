package io

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Stdin struct {
	ID
	*stdchannels
}

func NewStdin() (stdin *Stdin, err error) {
	stdin = new(Stdin)
	err = stdin.Init()
	return stdin, err
}

func (stdin *Stdin) Init() (err error) {
	if len(stdin.ID) == 0 {
        stdin.ID = ID("stdin://")
	} else {
        stdin.ID = ID("stdin://"+string(stdin.ID)+"@")
    }
	stdin.stdchannels = &stdchannels{}
	stdin.out = make(chan string)
	stdin.exit = make(chan struct{})
	go stdin.Handler()
	return nil
}

func (stdin *Stdin) Handler() {
	buffer := bufio.NewReader(os.Stdin)

LOOP:
	for {
		if stdin.IsDone() {
			break LOOP
		}
		out, err := buffer.ReadString('\n')
        fmt.Println("GOT", out)
		if err != nil {
			break LOOP //EOF
		}
		if stdin.IsDone() {
			break LOOP
		}
		// strings.TrimRight("","")
		stdin.out <- strings.TrimRight(out, "\n")
	}
	fmt.Println("finished stdin reader")
}
