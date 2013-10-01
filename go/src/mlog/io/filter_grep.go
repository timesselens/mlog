package io

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Grep struct {
	ID
	Re    string
	regex regexp.Regexp
	*stdchannels
}

func NewGrep(re string) (grep *Grep, err error) {
	grep = new(Grep)
    grep.Re = re
    err = grep.Init()
	return grep,err
}

func (grep *Grep) Init() (err error) {
    if grep.Re[0:1] == "/" && grep.Re[len(grep.Re)-1:] == "/" {
        grep.Re = strings.Trim(grep.Re,"/");
    }
	regex, err := regexp.Compile(grep.Re)
	if err != nil {
        grep.lastError = err.Error()
		return errors.New("unable to compile regex")
	}
    if len(grep.ID) == 0 {
        grep.ID = ID("grep://?re="+grep.Re)
    } else {
        grep.ID = ID("grep://"+string(grep.ID)+"@?re="+grep.Re)
    }
	grep.regex = *regex
    grep.stdchannels = &stdchannels{}
    grep.in = make(chan string)
    grep.out = make(chan string)
    grep.exit = make(chan struct{})
	go grep.Handler()
	return nil
}

func (grep *Grep) Handler() {
LOOP:
	for {
		select {
		case in, ok := <-grep.in:
            fmt.Printf("io.Grep got input %s matching with %s\n",in,grep.Re);
			if !ok {
                // grep.Err = errors.New("io.Grep: unable to read from input channel")
				break LOOP
			}
			if grep.regex.Match([]byte(in)) {
                fmt.Printf("grep MATCH %#v %s %s\n",grep,in, ok)
				grep.out <- "/" + grep.regex.String() + "/ " + in
			}
		case <-grep.exit:
			break LOOP
		}
	}
    fmt.Println("grep handler finished")
}
