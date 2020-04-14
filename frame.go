package cdp

import (
	"encoding/json"
	"github.com/4ydx/cdp/protocol/dom"
	"github.com/gorilla/websocket"
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

	// CurrentAction stores the Action that is currently active.
	CurrentAction *Action

	// LogLevel specifies how much information should be f.Browser.Logged. Higher number results in more data.
	LogLevel LogLevelValue
}

// SetCurrentAction sets the current action that the frame is evaluating.
func (f *Frame) SetCurrentAction(act *Action) {
	f.Lock()
	defer f.Unlock()
	f.CurrentAction = act
	f.ActionChan <- f.toJSON()
}

// SetDOM allows for setting the Frame DOM value safely.
func (f *Frame) SetDOM(dom *dom.GetFlattenedDocumentReply) {
	f.Lock()
	defer f.Unlock()
	f.DOM = dom
}

func (f *Frame) setChildNodes(nodes *[]dom.Node) {
	if nodes == nil {
		return
	}
	for _, node := range *nodes {
		if node.ChildNodeCount > 0 {
			f.setChildNodes(node.Children)
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

// Children returns a deep copy of the child nodes of the given parentID.
// NOTE: Expecting that code elsewhere has already populated the frame.DOM object.
func (f *Frame) Children(parentID dom.NodeID) []dom.Node {
	f.RLock()
	defer f.RUnlock()

	return f.children(parentID, []dom.Node{})
}

func (f *Frame) children(parentID dom.NodeID, found []dom.Node) []dom.Node {
	if f.DOM == nil {
		return found
	}
	for _, node := range f.DOM.Nodes {
		if node.ParentID != parentID {
			continue
		}
		found = append(found, node)

		if node.ChildNodeCount > 0 {
			// Remember that frame.DOM is a flattened representation so we should not be interacting with node.Children.
			found = f.children(node.NodeID, found)
		}
	}
	return found
}

// FindByAttribute will search the existing cached DOM for nodes whose given attribute matches the given value starting at the root specified by nodeID.
// NOTE: Expecting that code elsewhere has already populated the frame.DOM object.
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

// Stop closes used resources.
func (f *Frame) Stop(closeBrowser bool) {
	defer func() {
		err := f.Conn.Close()
		if err != nil {
			panic(err)
		}
		if closeBrowser && f.Browser != nil {
			f.Browser.Stop()
		}
	}()
	f.AllComplete <- struct{}{}
}

// IsCommandComplete indicates that all commands are complete.
func (f *Frame) IsCommandComplete() bool {
	f.RLock()
	defer f.RUnlock()

	return f.CurrentAction.CommandIndex == len(f.CurrentAction.Commands)
}

// IsComplete indicates that all commands and events are complete.
func (f *Frame) IsComplete() bool {
	f.RLock()
	defer f.RUnlock()

	complete := true
	for _, e := range f.CurrentAction.Events {
		if e.IsRequired && !e.IsFound {
			complete = false
		}
	}
	return f.CurrentAction.CommandIndex == len(f.CurrentAction.Commands) && complete
}

// CommandTimeout once timed out will trigger an error and stop the automation.
func (f *Frame) CommandTimeout() <-chan time.Time {
	f.RLock()
	defer f.RUnlock()

	return time.After(f.CurrentAction.Commands[f.CurrentAction.CommandIndex].Timeout)
}

// ToJSON encodes the current command.  This is the chrome devtools protocol request.
// In the event that all commands are complete, continue to display the last command for debugging convenience.
func (f *Frame) ToJSON() []byte {
	f.RLock()
	defer f.RUnlock()
	return f.toJSON()
}

func (f *Frame) toJSON() []byte {
	index := f.CurrentAction.CommandIndex
	if f.CurrentAction.CommandIndex == len(f.CurrentAction.Commands) {
		index--
	}
	s := f.CurrentAction.Commands[index]

	j, err := json.Marshal(s)
	if err != nil {
		f.Browser.Log.Fatal(err)
	}
	return j
}

// Log writes the current state of the action to the f.Browser.Log.
func (f *Frame) Log() {
	f.RLock()
	defer f.RUnlock()

	f.Browser.Log.Printf("Frame %+v\n", f.CurrentAction)
	for i, command := range f.CurrentAction.Commands {
		f.Browser.Log.Printf("%d Command %d Params %+v", i, command.ID, command.Params)
		f.Browser.Log.Printf("%d Command %d Return %+v", i, command.ID, command.Reply)
	}
}

// HasCommandID determines if an id matches the current action's command's unique id.
func (f *Frame) HasCommandID(id int64) bool {
	f.RLock()
	defer f.RUnlock()

	if f.CurrentAction.CommandIndex == len(f.CurrentAction.Commands) {
		return false
	}
	return f.CurrentAction.Commands[f.CurrentAction.CommandIndex].ID == id
}

// HasEvent returns true when the action has an event with the given MethodType.
func (f *Frame) HasEvent(name string) bool {
	f.RLock()
	defer f.RUnlock()

	_, ok := f.CurrentAction.Events[name]
	return ok
}

// GetCommandMethod returns the method of the command that is currently active or the very last method.
func (f *Frame) GetCommandMethod() string {
	f.RLock()
	defer f.RUnlock()

	if f.CurrentAction.CommandIndex == len(f.CurrentAction.Commands) {
		return f.CurrentAction.Commands[f.CurrentAction.CommandIndex-1].Method
	}
	return f.CurrentAction.Commands[f.CurrentAction.CommandIndex].Method
}

// SetEvent takes the given message and sets an event's params or results's.
func (f *Frame) SetEvent(frame *Frame, name string, m Message) error {
	f.Lock()
	defer f.Unlock()

	// Attempt to compare the incoming Event's frameID value with the existing value.
	if e, ok := f.CurrentAction.Events[name]; ok {
		if frame.FrameID == "" {
			f.Browser.Log.Println(".ERR FrameID is empty during event processing.")
			if len(m.Params) > 0 {
				err := e.Value.UnmarshalJSON(m.Params)
				if err != nil {
					f.Browser.Log.Printf("Unmarshal params error: %s; for %+v; from %+v", err.Error(), e.Value, m.Params)
					return err
				}
			} else {
				err := e.Value.UnmarshalJSON(m.Result)
				if err != nil {
					f.Browser.Log.Printf("Unmarshal result error: %s; for %+v; from %+v", err.Error(), e.Value, m.Result)
					return err
				}
			}
		} else {
			if len(m.Params) > 0 {
				if ok, err := e.Value.MatchFrameID(frame.FrameID, m.Params); !ok {
					if err != nil {
						f.Browser.Log.Printf("Unmarshal error: %s", err)
						return err
					}
					// When the frameID does not match, it is definitely not intended for the current Frame.
					f.Browser.Log.Printf("No matching frameID %s %s", m.Method, m.Params)
					return nil
				}
			} else {
				if ok, err := e.Value.MatchFrameID(frame.FrameID, m.Result); !ok {
					if err != nil {
						f.Browser.Log.Printf("Unmarshal error: %s", err)
						return err
					}
					f.Browser.Log.Printf("No matching frameID %s %s", m.Method, m.Result)
					return nil
				}
			}
		}
		UpdateDOMEvent(frame, m.Method, e.Value)

		e.IsFound = true
		f.CurrentAction.Events[string(name)] = e

		f.Browser.Log.Printf(".EVT: %s %+v\n", name, m)
		if frame.LogLevel >= LogDetails {
			f.Browser.Log.Printf("    : %+v\n", e)
			f.Browser.Log.Printf("    : %+v\n", e.Value)
		}
	}
	return nil
}

// SetResult applies the message returns to the current command and advances the command.
func (f *Frame) SetResult(frame *Frame, m Message) error {
	f.Lock()
	defer f.Unlock()

	s := f.CurrentAction.Commands[f.CurrentAction.CommandIndex]
	if frame.FrameID == "" {
		err := s.Reply.UnmarshalJSON(m.Result)
		if err != nil {
			f.Browser.Log.Printf("Unmarshal error: %s", err)
			return err
		}
		frame.FrameID = s.Reply.GetFrameID()
	} else {
		if ok, err := s.Reply.MatchFrameID(frame.FrameID, m.Result); !ok {
			if err != nil {
				f.Browser.Log.Printf("Unmarshal error: %s", err)
				return err
			} else {
				f.Browser.Log.Printf("No matching frameID")
				return nil
			}
		}
	}
	f.CurrentAction.CommandIndex++

	f.Browser.Log.Printf(".STP COMPLETE: %+v\n", s)
	if frame.LogLevel >= LogDetails {
		f.Browser.Log.Printf("             : %+v\n", s.Params)
		f.Browser.Log.Printf("             : %+v\n", s.Reply)
	}
	return nil
}

// Clear the action.
func (f *Frame) Clear() {
	f.Lock()
	defer f.Unlock()

	f.CurrentAction = &Action{}
}
