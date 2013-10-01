package core

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var Config struct {
    sync.RWMutex
    Verbose bool
    Webserver struct {
        HTTPListenAddress string
        WebSocketListenAddress string
    }
    Debug struct {
        Core bool
        HTTP bool
        Process bool
    }
}

// Global channels
var Channels struct {
	Exit chan bool
}

func init() {
	Channels.Exit = make(chan bool)
}

func Start() {
	reallyExit := false
	i := 4
	for {
		select {
		case v := <-Channels.Exit:
			i = 4
			reallyExit = v
            if reallyExit && (Config.Verbose || Config.Debug.Core) {
                fmt.Println("reveived true on Exit channel, exitting unless false is received within 5 sec")
            }
		case <-time.After(1 * time.Second):
			if reallyExit {
				i = i - 1
				if i == 0 {
					os.Exit(0)
				}
			}
		}
	}
}
