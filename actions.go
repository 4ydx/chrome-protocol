package cdp

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// Wait is the default timeout taken as the action wait loop runs.
// The loop's select is triggered and the action will be checked for completion.
var Wait = time.Millisecond * 50

// CommandReply specifies required methods for handling the json encoded replies received from the server.
type CommandReply interface {
	json.Unmarshaler
	MatchFrameID(frameID string, m []byte) bool
	GetFrameID() string
}

// Event holds the value returned by the server based on a matching MethodType name.
type Event struct {
	Name       string
	Value      CommandReply
	IsRequired bool
	IsFound    bool
}

// Command represents a single json request sent to the server over the websocket.
type Command struct {
	// Values required to make a chrome devtools protocol request.
	ID     int64          `json:"id"`
	Method string         `json:"method,omitempty"`
	Params json.Marshaler `json:"params,omitempty"`

	Reply           CommandReply  `json:"-"` // The struct that will be filled when a matching command Id is found in a reply over the chrome websocket.
	Timeout         time.Duration `json:"-"` // How long until the current command experiences a timeout, which will halt the entire process.
	PreviousReturns func()        `json:"-"` // Method defined by the user to take the previous command's Returns and apply them to the current command's Params.
}

// Action represents a collection of json requests (commands) and any events that those requests might trigger that need to be tracked.
type Action struct {
	*sync.RWMutex
	Commands     []Command
	CommandIndex int
	Events       map[string]Event
	Start        *time.Time
	Frame        *Frame
}

// NewAction returns a newly created action with any events that will be triggered by commands the action will take.
func NewAction(page *Frame, events []Event, commands []Command) *Action {
	act := &Action{
		RWMutex:  &sync.RWMutex{},
		Events:   make(map[string]Event),
		Commands: commands,
		Frame:    page,
	}
	for _, e := range events {
		act.Events[e.Name] = e
	}
	return act
}

// IsCommandComplete indicates that all commands are complete.
func (act *Action) IsCommandComplete() bool {
	act.RLock()
	defer act.RUnlock()

	return act.CommandIndex == len(act.Commands)
}

// IsComplete indicates that all commands and events are complete.
func (act *Action) IsComplete() bool {
	act.RLock()
	defer act.RUnlock()

	complete := true
	for _, e := range act.Events {
		if e.IsRequired && !e.IsFound {
			complete = false
		}
	}
	return act.CommandIndex == len(act.Commands) && complete
}

// CommandTimeout once timed out will trigger an error and stop the automation.
func (act *Action) CommandTimeout() <-chan time.Time {
	act.RLock()
	defer act.RUnlock()

	return time.After(act.Commands[act.CommandIndex].Timeout)
}

// ToJSON encodes the current command.  This is the chrome devtools protocol request.
// In the event that all commands are complete, continue to display the last command for debugging convenience.
func (act *Action) ToJSON() []byte {
	act.RLock()
	defer act.RUnlock()

	if act.Start == nil {
		t := time.Now()
		act.Start = &t
	}
	index := act.CommandIndex
	if act.CommandIndex == len(act.Commands) {
		index--
	}
	s := act.Commands[index]

	j, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return j
}

// Run sends the current action to websocket code that will create a request.
// Then the action will wait until all commands and expected events are completed.
func (act *Action) Run() error {
	Cache.Set(act)
	ActionChan <- act.ToJSON()
	commandTimeout := act.CommandTimeout()
	for {
		select {
		case <-commandTimeout:
			// The current action's current command has timed out.
			return fmt.Errorf("command timeout %s", act.ToJSON())
		case <-CacheCompleteChan:
			// The current action is complete.
			return nil
		case commandTimeout = <-CommandChan:
			// Set the current timeout to the next command's timeout.
			log.Print("Next command timeout set.")
		}
	}
}

// Log writes the current state of the action to the log.
func (act *Action) Log() {
	act.RLock()
	defer act.RUnlock()

	log.Printf("Action %+v\n", act)
	for i, command := range act.Commands {
		log.Printf("%d Command %d Params %+v", i, command.ID, command.Params)
		log.Printf("%d Command %d Return %+v", i, command.ID, command.Reply)
	}
}

// HasCommandID determines if an id matches the current action's command's unique id.
func (act *Action) HasCommandID(id int64) bool {
	act.RLock()
	defer act.RUnlock()

	if act.CommandIndex == len(act.Commands) {
		return false
	}
	return act.Commands[act.CommandIndex].ID == id
}

// HasEvent returns true when the action has an event with the given MethodType.
func (act *Action) HasEvent(name string) bool {
	act.RLock()
	defer act.RUnlock()

	_, ok := act.Events[name]
	return ok
}

// GetCommandMethod returns the method of the command that is currently active or the very last method.
func (act *Action) GetCommandMethod() string {
	act.RLock()
	defer act.RUnlock()

	if act.CommandIndex == len(act.Commands) {
		return act.Commands[act.CommandIndex-1].Method
	}
	return act.Commands[act.CommandIndex].Method
}

// GetFrameID returns the frameID of the current frame.
func (act *Action) GetFrameID() string {
	act.RLock()
	defer act.RUnlock()

	return act.Frame.FrameID
}

// SetEvent takes the given message and sets an event's params or results's.
func (act *Action) SetEvent(name string, m Message) error {
	act.Lock()
	defer act.Unlock()

	// Attempt to compare the incoming Event's frameID value with the existing value.
	frameID := act.Frame.GetFrameID()
	if e, ok := act.Events[name]; ok {
		if frameID == "" {
			log.Println(".ERR FrameID is empty during event processing.")
			if len(m.Params) > 0 {
				err := e.Value.UnmarshalJSON(m.Params)
				if err != nil {
					log.Printf("Unmarshal params error: %s; for %+v; from %+v", err.Error(), e.Value, m.Params)
					return err
				}
			} else {
				err := e.Value.UnmarshalJSON(m.Result)
				if err != nil {
					log.Printf("Unmarshal result error: %s; for %+v; from %+v", err.Error(), e.Value, m.Result)
					return err
				}
			}
		} else {
			if len(m.Params) > 0 {
				if ok := e.Value.MatchFrameID(frameID, m.Params); !ok {
					// When the frameID does not match, it is definitely not intended for the current Action.
					log.Printf("No matching frameID %s %s", m.Method, m.Params)
					return nil
				}
			} else {
				if ok := e.Value.MatchFrameID(frameID, m.Result); !ok {
					log.Printf("No matching frameID %s %s", m.Method, m.Result)
					return nil
				}
			}
		}
		e.IsFound = true
		act.Events[string(name)] = e

		log.Printf(".EVT: %s %+v\n", name, m)
		if LogLevel >= LogDetails {
			log.Printf("    : %+v\n", e)
			log.Printf("    : %+v\n", e.Value)
		}
	}
	return nil
}

// SetResult applies the message returns to the current command and advances the command.
func (act *Action) SetResult(m Message) error {
	act.Lock()
	defer act.Unlock()

	s := act.Commands[act.CommandIndex]
	frameID := act.Frame.GetFrameID()
	if frameID == "" {
		err := s.Reply.UnmarshalJSON(m.Result)
		if err != nil {
			log.Printf("Unmarshal error: %s", err)
			return err
		}
		act.Frame.SetFrameID(s.Reply.GetFrameID())
	} else {
		if ok := s.Reply.MatchFrameID(frameID, m.Result); !ok {
			log.Printf("No matching frameID")
			return nil
		}
	}
	act.CommandIndex++

	log.Printf(".STP COMPLETE: %+v\n", s)
	if LogLevel >= LogDetails {
		log.Printf("             : %+v\n", s.Params)
		log.Printf("             : %+v\n", s.Reply)
	}
	return nil
}
