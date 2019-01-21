package cdp

import (
	"encoding/json"
	"fmt"
	"time"
)

// Wait is the default timeout taken as the action wait loop runs.
// The loop's select is triggered and the action will be checked for completion.
var Wait = time.Millisecond * 50

// CommandReply specifies required methods for handling the json encoded replies received from the server.
type CommandReply interface {
	json.Unmarshaler
	MatchFrameID(frameID string, m []byte) (bool, error)
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

	Reply   CommandReply  `json:"-"` // The struct that will be filled when a matching command Id is found in a reply over the chrome websocket.
	Timeout time.Duration `json:"-"` // How long until the current command experiences a timeout, which will halt the entire process.
}

// Action represents a collection of json requests (commands) and any events that those requests might trigger that need to be tracked.
type Action struct {
	Commands     []Command
	CommandIndex int
	Events       map[string]Event
}

// NewAction returns a newly created action with any events that will be triggered by commands the action will take.
func NewAction(events []Event, commands []Command) *Action {
	act := &Action{
		Events:   make(map[string]Event),
		Commands: commands,
	}
	for _, e := range events {
		act.Events[e.Name] = e
	}
	return act
}

// Run sends the current action to websocket code that will create a request.
// Then the action will wait until all commands and expected events are completed.
func (act *Action) Run(frame *Frame) error {
	frame.SetCurrentAction(act)
	commandTimeout := frame.CommandTimeout()
	for {
		select {
		case <-commandTimeout:
			// The current action's current command has timed out.
			return fmt.Errorf("command timeout %s", frame.ToJSON())
		case <-frame.CacheCompleteChan:
			// The current action is complete.
			return nil
		case commandTimeout = <-frame.CommandChan:
			// Set the current timeout to the next command's timeout.
			frame.Browser.Log.Print("Next command timeout set.")
		}
	}
}
