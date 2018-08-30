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

// SetDOM allows for setting the Frame DOM value safely.
func (f *Frame) SetDOM(dom *dom.GetFlattenedDocumentReply) {
	f.Lock()
	defer f.Unlock()
	f.DOM = dom
}

// SetChildNodes updates the Frame DOM with new child nodes.
func (f *Frame) SetChildNodes(nodes *[]dom.Node) {
	f.Lock()
	defer f.Unlock()
	f.setChildNodesHelper(nodes)
}

func (f *Frame) setChildNodesHelper(nodes *[]dom.Node) {
	if nodes == nil {
		return
	}
	for _, node := range *nodes {
		if node.ChildNodeCount > 0 {
			f.setChildNodesHelper(node.Children)
		}
		f.DOM.Nodes = append(f.DOM.Nodes, node)
	}
}

// GetDOM allows for getting the Frame DOM value safely.
// This could be a bit racy depending on when documentUpdated events are fired.
func (f *Frame) GetDOM() *dom.GetFlattenedDocumentReply {
	f.RLock()
	defer f.RUnlock()
	return f.DOM
}

// AddDOMNode allows for setting the Frame DOM value safely.
func (f *Frame) AddDOMNode(node dom.Node) {
	f.Lock()
	defer f.Unlock()
	f.DOM.Nodes = append(f.DOM.Nodes, node)
}

// FindByAttribute will search the existing cached DOM for nodes whose given attribute matches the given value.
// Note that this will not attempt get the DOM if it isn't yet populated.  Might change this later.
func (f *Frame) FindByAttribute(parentID dom.NodeID, attribute, value string) []dom.Node {
	f.RLock()
	defer f.RUnlock()

	return f.findByAttributeHelper(parentID, attribute, value, []dom.Node{})
}

func (f *Frame) findByAttributeHelper(parentID dom.NodeID, attribute, value string, found []dom.Node) []dom.Node {
	if f.DOM == nil {
		return found
	}
	for _, node := range f.DOM.Nodes {
		if node.ParentID != parentID {
			continue
		}
		if node.Attributes != nil {
			for i, attr := range *node.Attributes {
				if attr == attribute && (*node.Attributes)[i+1] == value {
					found = append(found, node)
				}
			}
		}
		if node.ChildNodeCount > 0 {
			// Remember that frame.DOM is a flattened representation so we should not be interacting with node.Children.
			found = f.findByAttributeHelper(node.NodeID, attribute, value, found)
		}
	}
	return found
}

// GetFrameID returns the current frameID.
func (f *Frame) GetFrameID() string {
	f.RLock()
	defer f.RUnlock()

	return f.FrameID
}

// SetFrameID sets the current frameID.
func (f *Frame) SetFrameID(frameID string) {
	f.Lock()
	defer f.Unlock()

	f.FrameID = frameID
}

// Stop closes used resources.
func (f *Frame) Stop(closeBrowser bool) {
	defer func() {
		err := f.Conn.Close()
		if err != nil {
			panic(err)
		}
		if closeBrowser {
			f.Browser.Stop()
		}
	}()
	f.AllComplete <- struct{}{}
}
