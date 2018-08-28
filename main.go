package cdp

import (
	"github.com/gorilla/websocket"
	"log"
	"os"
	"sync"
	"time"
)

var (
	// Conn is the connection to the websocket.
	Conn *websocket.Conn

	// AllComplete will trigger a close on the websocket.
	// Typically AllComplete or the OsInterrupt channels will fire and the write loop will send a request to close the socket.
	AllComplete chan struct{}

	// CacheCompleteChan sends the signal that the cached action is completed (all commands and events).
	CacheCompleteChan chan struct{}

	// CommandChan sends the signal that a command has been completed and an Action can advance.
	CommandChan chan (<-chan time.Time)

	// ActionChan sends Actions to the websocket.
	ActionChan chan []byte

	// Cache stores the Action that is currently active.
	Cache *ActionCache

	// LogFile is the file that all log output will be written to.
	LogFile *os.File

	// LogLevel specifies how much information should be logged. Higher number results in more data.
	LogLevel LogLevelValue
)

// LogLevelValue is the type for loglevel information.
type LogLevelValue int

const (
	// LogBasic records outgoing commands, their replies, and any specified events.
	LogBasic = LogLevelValue(0)
	// LogDetails records additional details about the reply from the server for a given command/event.
	LogDetails = LogLevelValue(1)
	// LogAll records everything.
	LogAll = LogLevelValue(2)
)

func init() {
	var err error
	LogFile, err = os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	//log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(LogFile)
}

// StartWithLog prepares required resources to begin automation and sets the std logger file.
func StartWithLog(port int, logFile string, logLevel LogLevelValue) *Frame {
	previous := LogFile
	defer previous.Close()

	var err error
	LogFile, err = os.Create(logFile)
	if err != nil {
		panic(err)
	}
	//log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(LogFile)
	return Start(port, logLevel)
}

// Start prepares required resources to begin automation.
func Start(port int, logLevel LogLevelValue) *Frame {
	Conn = GetWebsocket(port)
	Cache = &ActionCache{
		RWMutex: &sync.RWMutex{},
	}
	AllComplete = make(chan struct{})

	CacheCompleteChan = make(chan struct{})
	ActionChan = make(chan []byte)
	CommandChan = make(chan (<-chan time.Time))
	LogLevel = logLevel

	go Write(Conn, ActionChan, AllComplete)
	go Read(Conn, CommandChan, CacheCompleteChan, Cache)

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
