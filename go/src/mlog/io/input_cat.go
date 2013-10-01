package io
/*

import (
    "bufio"
    "fmt"
    "log"
	"os"
	"syscall"
    "mlog/core"
)
type Cat struct {
    core.Process
    file string
}

func NewCat(file string) (cat *Cat, err error) {
	cat = &Cat{core.Process{In: make(chan string), Out: make(chan string)},file}
    go cat.Handler()
	return cat, err
}

func (cat Cat) Handler() {
        fh, err := os.OpenFile(cat.file, syscall.O_RDONLY, 0666)
        defer func() {
            fh.Close()
            close(cat.In)
            close(cat.Out)
        }()

        if err != nil {
            log.Fatal(err)
        }
        scanner := bufio.NewScanner(fh)
        for scanner.Scan() {
            cat.Out <- scanner.Text()
        }
        if err := scanner.Err(); err != nil {
            fmt.Fprintln(os.Stderr, "error scanning file", err)
        }
    }
*/
