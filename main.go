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

	// AllComplete will trigger a close on the websocket.
	// Typically AllComplete or the OsInterrupt channels will fire and the write loop will send a request to close the socket.
	AllComplete chan struct{}

	// StepChan sends the signal that a step has been completed and an Action can advance.
	StepChan chan struct{}

	// ActionChan sends Actions to the websocket.
	ActionChan chan []byte

	// Cache stores the Action that is currently active.
	Cache *ActionCache

	// LogFile is the file that all log output will be written to.
	LogFile *os.File
)

func init() {
	LogFile, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	//log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(LogFile)
}

// Start prepares required resources to begin automation.
func Start(port int) *Frame {
	Conn = GetWebsocket(port)
	Cache = &ActionCache{
		RWMutex: &sync.RWMutex{},
	}
	AllComplete = make(chan struct{})

	ActionChan = make(chan []byte)
	StepChan = make(chan struct{})

	go Write(Conn, ActionChan, AllComplete)
	go Read(Conn, StepChan, Cache)

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
	defer func() {
		Conn.Close()
		LogFile.Close()
	}()
	AllComplete <- struct{}{}
}
