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
	StepChan    chan struct{}
	ActionChan  chan *Action
	Cache       *ActionCache
)

// Start prepares required resources to begin automation.
func Start() *ID {
	f, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(f)

	Conn = GetWebsocket()
	Cache = NewActionCache()
	ShutDown = make(chan struct{})
	AllComplete = make(chan struct{})

	ActionChan = make(chan *Action)
	StepChan = make(chan struct{})

	go Write(Conn, ActionChan, Cache, ShutDown, AllComplete)
	go Read(Conn, StepChan, Cache, ShutDown)

	id := &ID{
		RWMutex: &sync.RWMutex{},
		Value:   11111,
	}
	return id
}

// Stop closes used resources.
func Stop() {
	defer Conn.Close()

	AllComplete <- struct{}{}
	<-ShutDown
}
