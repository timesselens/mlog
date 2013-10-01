package control

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"syscall"
)

type NamedPipe struct {
	path string
	fh   *os.File
	in   chan string
	out  chan string
}

type PipeList struct {
	sync.RWMutex
	List map[string]*NamedPipe
}

type PipeServer struct {
    Pipes *PipeList
    New chan string
}

var Pipes = &PipeList{List: make(map[string]*NamedPipe)}

func NamedPipeServer(dir string) (p *PipeServer){

	for _, fifoname := range []string{"error", "access"} {
		f, err := createFifo(dir, fifoname)
		if err != nil {
			log.Fatal(err)
		}

		Pipes.Lock()
		Pipes.List[fifoname] = f
		Pipes.Unlock()
		// f.in <- fmt.Sprintf("starting controlserver in '%s'", dir)
	}

	for _, p := range Pipes.List {
		if p.fh == nil {
			Pipes.Lock()
			fh, err := os.OpenFile(p.path, syscall.O_RDWR, 0666)
			if err != nil {
				log.Fatal(err)
			}
            err = syscall.SetNonblock(int(fh.Fd()), true)
            if err != nil {
                log.Fatal("unable to set nonblocking mode on fifo buffer")
            }
			p.fh = fh
			fmt.Println(p)
			Pipes.Unlock()
		}
		go handlepipe(p)
	}

    return &PipeServer{Pipes: Pipes, New:make(chan string)}

}


func handlepipe(p *NamedPipe) {
    defer p.fh.Close()
    LOOP:
    for {
        select {
        case m := <-p.in:
            _, err := p.fh.WriteString(m + "\n")
            if err != nil {
                fmt.Println("| err",err)
                break LOOP;
            }
        }
    }
    close(p.in)
    close(p.out)
    Pipes.Lock()
    delete(Pipes.List,path.Base(p.path))
    Pipes.Unlock()
    fmt.Println(">>>>>>>>>>>>>>>>>>>> done handling pipe",p.path)
}

func createFifo(dir, name string) (fifo *NamedPipe, err error) {
	_, err = os.Stat(dir + "/" + name)

	if os.IsNotExist(err) {
		syscall.Mknod(dir+"/"+name, syscall.S_IFIFO|0666, 0)
		err = nil
	}

	n := &NamedPipe{path: dir + "/" + name, in: make(chan string), out: make(chan string)}

	return n, err
}

func (p PipeList) Send(np, m string) (err error) {

	pipe,ok := p.List[np]
	if !ok {
		return errors.New("named pipe '" + np + "' not found")
	}

	go func() {
		pipe.in <- m
	}()

	return err
}
