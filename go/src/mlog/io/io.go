package io

import (
	"errors"
	"fmt"
)

type ID string
type IDer interface {
	GetID() (id string)
	SetID(id string)
}

func (src *ID) GetID() string {
	return string(*src)
}

func (src *ID) SetID(id string) {
	p := ID(id)
	src = &p
}
func (src *ID) String() string {
	return string(*src)
}

type stdchannels struct {
	in   chan string
	out  chan string
	err  chan string
	exit chan struct{}
    pipe_from *Processor
    pipe_to *Processor
	sub  map[string][]chan string
    lastError string
}

type Processor interface {
	GetID() (id string)
	SetID(id string)
	Init() (err error)
	Handler()
	Send(s string)
	Recv() (s string, ok bool)
    GetSourceProcessor() (s *Processor)
    GetTargetProcessor() (t *Processor)
    IsPiped() bool
	Close()
    IsDone() (bool)
	Exit()
	pipe(t Processor)
    setPipeSrc(t Processor)
    LastError() string
}

func Pipe(first Processor, second Processor, rest ...Processor) (out Processor) {
	list := []Processor{first, second}
	list = append(list, rest...)
	for i, p := range list {
		if i <= len(list)-2 {
			// fmt.Printf("connecting [%d]%T to [%d]%T\n",i,p,i+1,list[i+1])
            if list[i+1] != nil {
                p.pipe(list[i+1])
                list[i+1].setPipeSrc(p)
            }
		}
	}
	return list[0]
}

// func (p Processor) Foo(out Processor) {
// }

func (c *stdchannels) pipe(out Processor) {
    c.pipe_to = &out
	go func() {
	LOOP:
		for {
			select {
			case val, ok := <-c.out:
				if !ok {
					fmt.Println("Error receiving from channel")
					break LOOP
				}
				if out.IsDone() {
                    fmt.Println("%T is done breaking loop",out)
					break LOOP
				}
                fmt.Println("about to pipe",val)
				out.Send(val)
			case <-c.exit:
				break LOOP
			}

		}
		fmt.Printf("broken pipe %T %T\n", c, out)
	}()
	return
}

func (c *stdchannels) Subscribe(stdtype string, s chan string) (err error) {
	if stdtype != "out" || stdtype != "err" || stdtype != "in" {
		return errors.New("subscription to unknown channel")
	}
	if c.sub == nil {
		c.sub = make(map[string][]chan string)
	}
	c.sub[stdtype] = append(c.sub[stdtype], s)
	return nil
}

// func (c *stdchannels) ID() (id string) {
    // return fmt.Sprintf("%T",c)
// }

func (c *stdchannels) Send(s string) {
	if c.in != nil {
		c.in <- s
	}
}

func (c *stdchannels) Recv() (s string, ok bool) {
	s, ok = <-c.out
	return s, ok
}

// Set the exit channel to nil so the select will unblock on <-type.exit
func (c *stdchannels) Exit() {
	c.exit = nil
}

// Drains and closes stdin, stdout and stderr channels
func (c *stdchannels) Close() {
	if c.in != nil {
		select {
		case <-c.in:
		default:
		}
		close(c.in)
	}
	if c.out != nil {
		select {
		case <-c.out:
		default:
		}
		close(c.out)
	}
	if c.err != nil {
		select {
		case <-c.err:
		default:
		}
		close(c.err)
	}
}

func (c *stdchannels) IsDone() bool {
    return c.exit == nil
}

func (c *stdchannels) IsPiped() bool {
    return c.pipe_from != nil || c.pipe_to != nil
}

func (c *stdchannels) LastError() string {
    return c.lastError
}

func (c *stdchannels) GetSourceProcessor() (t *Processor) {
    return c.pipe_from
}

func (c *stdchannels) GetTargetProcessor() (t *Processor) {
    return c.pipe_to
}

func (c *stdchannels) setPipeSrc(t Processor) {
    c.pipe_from = &t
}
