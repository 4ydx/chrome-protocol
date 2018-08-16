package cdp

import (
	"github.com/gorilla/websocket"
	"log"
	"os"
	"sync"
)

var (
	Conn        *websocket.Conn
	AllComplete chan struct{}
	ShutDown    chan struct{}
)

// Start prepares required resources to begin automation.
func Start() (*ActionCache, *ID, chan *Action, chan bool) {
	f, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(f)

	Conn = GetWebsocket()
	ShutDown = make(chan struct{})
	AllComplete = make(chan struct{})

	actionChan := make(chan *Action)
	stepChan := make(chan bool)

	actionCache := NewActionCache()

	go Write(Conn, actionChan, actionCache, ShutDown, AllComplete)
	go Read(Conn, stepChan, actionCache, ShutDown)

	id := &ID{
		RWMutex: &sync.RWMutex{},
		Value:   11111,
	}
	return actionCache, id, actionChan, stepChan
}

// Stop closes used resources.
func Stop() {
	defer Conn.Close()

	AllComplete <- struct{}{}
	<-ShutDown
}
