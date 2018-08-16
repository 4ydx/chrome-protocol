package cdp

import (
	"github.com/gorilla/websocket"
	"log"
	"os"
	"sync"
)

var Conn *websocket.Conn

// Start prepares required resources to begin automation.
func Start() (*EventCache, *ID, chan *Action, chan bool, chan struct{}, chan struct{}) {
	f, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(f)

	Conn = GetWebsocket()

	actionChan := make(chan *Action)
	stepChan := make(chan bool)

	shutdown := make(chan struct{})
	allComplete := make(chan struct{})

	actionCache := NewActionCache()
	eventCache := NewEventCache()

	go Write(Conn, actionChan, actionCache, shutdown, allComplete)
	go Read(Conn, stepChan, actionCache, eventCache, shutdown)

	id := &ID{
		RWMutex: &sync.RWMutex{},
		Value:   11111,
	}
	return eventCache, id, actionChan, stepChan, allComplete, shutdown
}

// Stop closes used resources.
func Stop(allComplete chan<- struct{}, shutdown <-chan struct{}) {
	defer Conn.Close()

	allComplete <- struct{}{}
	<-shutdown
}
