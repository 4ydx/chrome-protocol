package cdp

import (
	"github.com/gorilla/websocket"
	"log"
	"os"
	"sync"
)

var (
	// Conn is the connection to the websocket.
	Conn *websocket.Conn

	// AllComplete receives empty structs and will trigger a close on the websocket.
	// Typically AllComplete or the OsInterrupt channels will fire, triggering a close
	// of the websocket connection via the write loop.  This, in turn, will cause the
	// the write loop to wait for the shutdown channel to be closed or a timeout to be issued.
	AllComplete chan struct{}

	// ShutDown will be closed when reading the websocket is no longer possible.
	ShutDown chan struct{}

	// StepChan sends the signal that a step has been completed and an Action can advance.
	StepChan chan struct{}

	// ActionChan sends Actions to the websocket.
	ActionChan chan *Action

	// Cache stores the Action that is currently active.
	Cache *ActionCache
)

// Start prepares required resources to begin automation.
func Start(port int) *Frame {
	f, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	//log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(f)

	Conn = GetWebsocket(port)
	Cache = &ActionCache{}
	ShutDown = make(chan struct{})
	AllComplete = make(chan struct{})

	ActionChan = make(chan *Action)
	StepChan = make(chan struct{})

	go Write(Conn, ActionChan, Cache, ShutDown, AllComplete)
	go Read(Conn, StepChan, Cache, ShutDown)

	page := &Frame{
		RWMutex: &sync.RWMutex{},
		RequestID: RequestID{
			RWMutex: &sync.RWMutex{},
			Value:   11111,
		},
	}
	return page
}

// Stop closes used resources.
func Stop() {
	defer Conn.Close()

	AllComplete <- struct{}{}
}
