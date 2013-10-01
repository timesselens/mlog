package control

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	. "util"
)

type UnixDomainSocket struct {
	path   string
	handle func(c net.Conn)
}

type UnixDomainSocketServer struct {
	sockets []*UnixDomainSocket
}

func UnixDomainServer(dir string) {
	echo := &UnixDomainSocket{path: dir+"/echo", handle: handleEcho}
	mem := &UnixDomainSocket{path: dir+"/mem", handle: handleMem}
	memH := &UnixDomainSocket{path: dir+"/memory", handle: handleMemHuman}
    udss :=  &UnixDomainSocketServer{sockets: []*UnixDomainSocket{echo, mem, memH}}
	udss.Run()
}

func (u *UnixDomainSocketServer) Run() {
	for _, socket := range u.sockets {
		go socket.Listen()
	}
}

func (u *UnixDomainSocket) Listen() (err error) {
	os.Remove(u.path)
	l, err := net.Listen("unix", u.path)
	fmt.Println("listening on", u.path)
	if err != nil {
		log.Fatal("error listening", err)
	}
	for {
		fd, err := l.Accept()
		if err != nil {
			println("accept error", err)
			return err
		}

		go u.handle(fd)
	}

}

func handleEcho(c net.Conn) {
	defer c.Close()
	for {
		input, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Print("error", err)
			return
		}

		_, err = c.Write([]byte(input))
		if err != nil {
			log.Fatal("write error", err)
		}
	}
}

func handleMem(c net.Conn) {
	defer c.Close()

	MEM := new(runtime.MemStats)
	runtime.ReadMemStats(MEM)
	t := "#go %d\tmem: alloc: %d\t#mallocs: %d\t#frees: %d\t#diff: %d\tpause: %d\theap: %d\tstack: %d\n"
	r := fmt.Sprintf(t, runtime.NumGoroutine(), MEM.Alloc, MEM.Mallocs, MEM.Frees, MEM.Mallocs-MEM.Frees, MEM.PauseTotalNs, MEM.HeapAlloc, MEM.StackInuse)
	_, err := c.Write([]byte(r))
	if err != nil {
		log.Fatal("write error", err)
	}
}

func handleMemHuman(c net.Conn) {
	defer c.Close()

	MEM := new(runtime.MemStats)
	runtime.ReadMemStats(MEM)
	t := "#go %d\tmem: alloc: %s\t#mallocs: %d\t#frees: %d\t#diff: %d\tpause: %d\theap: %s\tstack: %s\n"
	r := fmt.Sprintf(t, runtime.NumGoroutine(), ByteSize(MEM.Alloc), MEM.Mallocs, MEM.Frees, MEM.Mallocs-MEM.Frees, MEM.PauseTotalNs, ByteSize(MEM.HeapAlloc), ByteSize(MEM.StackInuse))
	_, err := c.Write([]byte(r))
	if err != nil {
		log.Fatal("write error", err)
	}
}

// func (u *UnixDomainSocketServer) Add() {
// }
