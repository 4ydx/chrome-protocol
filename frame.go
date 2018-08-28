package cdp

import (
	"github.com/4ydx/cdp/protocol/dom"
	"github.com/gorilla/websocket"
	"os"
	"sync"
	"time"
)

// Frame stores the current FrameID.
type Frame struct {
	*sync.RWMutex
	DOM       *dom.GetFlattenedDocumentReply
	FrameID   string
	LoaderID  string
	RequestID RequestID

	Browser *Browser

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
}

// GetFrameID returns the current frameID.
func (p *Frame) GetFrameID() string {
	p.RLock()
	defer p.RUnlock()

	return p.FrameID
}

// SetFrameID sets the current frameID.
func (p *Frame) SetFrameID(frameID string) {
	p.Lock()
	defer p.Unlock()

	p.FrameID = frameID
}

// Stop closes used resources.
func (f *Frame) Stop(closeBrowser bool) {
	defer func() {
		f.Conn.Close()
		if closeBrowser {
			f.Browser.Stop()
		}
	}()
	f.AllComplete <- struct{}{}
}
