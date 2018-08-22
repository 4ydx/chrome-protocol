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

// StepReply specifies required methods for handling the json encoded replies received from the server.
type StepReply interface {
	json.Unmarshaler
	MatchFrameID(frameID string, m []byte) bool
	GetFrameID() string
}

// Event holds the value returned by the server based on a matching MethodType name.
type Event struct {
	Name       string
	Value      StepReply
	IsRequired bool
	IsFound    bool
}

// Step represents a single json request sent to the server over the websocket.
type Step struct {
	// Values required to make a chrome devtools protocol request.
	ID     int64          `json:"id"`
	Method string         `json:"method,omitempty"`
	Params json.Marshaler `json:"params,omitempty"`

	Reply           StepReply     `json:"-"` // The struct that will be filled when a matching step Id is found in a reply over the chrome websocket.
	Timeout         time.Duration `json:"-"` // How long until the current step experiences a timeout, which will halt the entire process.
	PreviousReturns func()        `json:"-"` // Method defined by the user to take the previous step's Returns and apply them to the current step's Params.
}

// Action represents a collection of json requests (steps) and any events that those requests might trigger that need to be tracked.
type Action struct {
	*sync.RWMutex
	Steps     []Step
	StepIndex int
	Events    map[string]Event
	Start     *time.Time
	Frame     *Frame
}

// NewAction returns a newly created action with any events that will be triggered by steps the action will take.
func NewAction(page *Frame, events []Event, steps []Step) *Action {
	act := &Action{
		RWMutex: &sync.RWMutex{},
		Events:  make(map[string]Event),
		Steps:   steps,
		Frame:   page,
	}
	for _, e := range events {
		act.Events[e.Name] = e
	}
	return act
}

// IsComplete indicates that all steps are complete.
func (act *Action) IsComplete() bool {
	act.RLock()
	defer act.RUnlock()

	return act.StepIndex == len(act.Steps)
}

// StepTimeout once timed out will trigger an error and stop the automation.
func (act *Action) StepTimeout() <-chan time.Time {
	act.RLock()
	defer act.RUnlock()

	s := act.Steps[act.StepIndex]
	return time.After(s.Timeout)
}

// ToJSON encodes the current step.  This is the chrome devtools protocol request.
// In the event that all steps are complete, continue to display the last step for debugging convenience.
func (act *Action) ToJSON() []byte {
	act.RLock()
	defer act.RUnlock()

	if act.Start == nil {
		t := time.Now()
		act.Start = &t
	}
	index := act.StepIndex
	if act.StepIndex == len(act.Steps) {
		index--
	}
	s := act.Steps[index]

	j, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return j
}

// Run sends the current action to websocket code that will create a request.
// Then the action will wait until all steps and expected events are completed.
func (act *Action) Run() error {
	Cache.Set(act)
	ActionChan <- act.ToJSON()
	stepTimeout := act.StepTimeout()
	for {
		select {
		case <-stepTimeout:
			return fmt.Errorf("step timeout %s", act.ToJSON())
		case <-StepChan:
			// Once a step is complete check to see if the action is entirely complete
			// or if it needs to continue on to the next step.
			if !act.IsComplete() {
				// Execute the next Step in the Action.
				ActionChan <- act.ToJSON()
				stepTimeout = act.StepTimeout()
			}
			if act.IsComplete() && Cache.EventsComplete() {
				log.Printf("Action Step completed %s %s", Cache.GetStepMethod(), Cache.GetFrameID())
				Cache.Clear()
				return nil
			}
			log.Printf("Action waiting %s %s", Cache.GetStepMethod(), Cache.GetFrameID())
		}
	}
}

// Log writes the current state of the action to the log.
func (act *Action) Log() {
	act.RLock()
	defer act.RUnlock()

	log.Printf("Action %+v\n", act)
	for i, step := range act.Steps {
		log.Printf("%d Step %d Params %+v", i, step.ID, step.Params)
		log.Printf("%d Step %d Return %+v", i, step.ID, step.Reply)
	}
}
