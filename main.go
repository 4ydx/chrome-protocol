package cdp

import (
	"sync"
	"time"
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

// Start prepares required resources to begin automation.
func Start(browser *Browser, logLevel LogLevelValue) *Frame {
	frame := &Frame{
		RWMutex: &sync.RWMutex{},
		RequestID: RequestID{
			RWMutex: &sync.RWMutex{},
			Value:   11111,
		},
		Browser:           browser,
		Conn:              GetWebsocket(browser.Port),
		Cache:             &ActionCache{RWMutex: &sync.RWMutex{}},
		CacheCompleteChan: make(chan struct{}),
		ActionChan:        make(chan []byte),
		CommandChan:       make(chan (<-chan time.Time)),
		AllComplete:       make(chan struct{}),
		LogLevel:          logLevel,
	}
	go Write(frame)
	go Read(frame)

	return frame
}
