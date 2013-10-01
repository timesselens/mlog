package io

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type WebSocket struct {
	ID
	Hostport string
	*stdchannels
	pool *webSocketClientPool
}

type webSocketClientPool struct {
	sync.RWMutex
	m map[*websocket.Conn]chan string
}

func NewWebSocket(hostport string) (ws *WebSocket, err error) {
	ws = new(WebSocket)
	ws.Hostport = hostport
	err = ws.Init()
	return ws, err
}

func (ws *WebSocket) Init() (err error) {
	if len(ws.Hostport) == 0 {
		ws.Hostport = "localhost:12345"
	}
	if len(ws.ID) == 0 {
		ws.ID = ID("ws://" + ws.Hostport + "/")
	} else {
		ws.ID = ID("ws://" + string(ws.ID) + "@" + ws.Hostport + "/")
	}
	ws.stdchannels = &stdchannels{}
	ws.in = make(chan string)
	ws.out = make(chan string)
	ws.err = make(chan string)
	ws.exit = make(chan struct{})
	ws.pool = &webSocketClientPool{m: make(map[*websocket.Conn]chan string)}
	go ws.ServerHandler()
	go ws.Handler()
	return nil
}

func (ws *WebSocket) ServerHandler() {
	server := http.Server{Addr: ws.Hostport, Handler: websocket.Handler(ws.socketHandler)}
	server.Handler = websocket.Handler(ws.socketHandler)
    err := server.ListenAndServe()
    if err != nil {
    log.Printf("unable to start websocket handler")
    }
}

func (ws *WebSocket) Handler() {
LOOP:
	for {
		select {
		case m, ok := <-ws.in:
			if !ok {
				log.Println("unable to read from in channel...")
				break LOOP
			}
			ws.pool.RLock()
			for _, bc := range ws.pool.m {
				bc <- m
			}
			ws.pool.RUnlock()
		case <-ws.exit:
			break LOOP
		}
	}
	log.Println("end of websocket handler")
}

func (ws *WebSocket) socketHandler(conn *websocket.Conn) {
	recvr := make(chan string)
	ws.pool.Lock()
	ws.pool.m[conn] = recvr
	ws.pool.Unlock()
	defer func() {
		close(recvr)
		ws.pool.Lock()
		delete(ws.pool.m, conn)
		ws.pool.Unlock()
	}()

LOOP:
	for {
		select {
		case m := <-recvr:
			fmt.Printf("websocket for client %s received %s\n", ws.Hostport, m)
			err := websocket.Message.Send(conn, m)
			if err != nil {
				log.Println("unable to send to socket... ")
				break LOOP
			}
		case <-ws.exit:
			break LOOP
		}
	}
	log.Println("end of socketHandler")
}
