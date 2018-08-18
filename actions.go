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

// Event holds the value returned by the server based on a matching MethodType name.
type Event struct {
	Name       string
	Value      json.Unmarshaler
	IsRequired bool
	IsFound    bool
}

// Step represents a single json request sent to the server over the websocket.
type Step struct {
	// Values required to make a chrome devtools protocol request.
	ID     int64          `json:"id"`
	Method string         `json:"method,omitempty"`
	Params json.Marshaler `json:"params,omitempty"`

	Reply           json.Unmarshaler `json:"-"` // The struct that will be filled when a matching step Id is found in a reply over the chrome websocket.
	Timeout         time.Duration    `json:"-"` // How long until the current step experiences a timeout, which will halt the entire process.
	PreviousReturns func()           `json:"-"` // Method defined by the user to take the previous step's Returns and apply them to the current step's Params.
}

// Action represents a collection of json requests (steps) and any events that those requests might trigger that need to be tracked.
// TODO: I'm not sure that multiple steps are really required here.  It might be less confusing if the Step was folded into the Action,
//       so that one action == one json api call across the websocket.
type Action struct {
	*sync.RWMutex
	Steps     []Step
	StepIndex int
	Events    map[string]Event
	Start     *time.Time
	Page      *Page
}

// NewAction returns a newly created action with any events that will be triggered by steps the action will take.
func NewAction(page *Page, events []Event, steps []Step) *Action {
	act := &Action{
		RWMutex: &sync.RWMutex{},
		Events:  make(map[string]Event),
		Steps:   steps,
		Page:    page,
	}
	for _, e := range events {
		act.Events[e.Name] = e
	}
	return act
}

// wait continues to query the state of the action.
// Once the action is complete, wait will return.
func (act *Action) wait(actionChan chan<- *Action, ac *ActionCache, stepChan <-chan struct{}) error {
	for {
		select {
		case <-time.After(Wait):
			if !act.IsComplete() && act.StepTimeout() {
				return fmt.Errorf("step timeout %s", act.ToJSON())
			}
			if act.IsComplete() && ac.EventsComplete() {
				log.Print("Action completed.")
				return nil
			}
			log.Print("Action waiting...")
		case <-stepChan:
			if !act.IsComplete() {
				if act.StepTimeout() {
					return fmt.Errorf("step timeout %s", act.ToJSON())
				}
				// Push the current action's next step to the server.
				actionChan <- act
			}
			if act.IsComplete() && ac.EventsComplete() {
				log.Printf("Action completed.")
				return nil
			}
			log.Printf("Action waiting...")
		}
	}
}

// IsComplete indicates that all steps are complete.
func (act *Action) IsComplete() bool {
	act.RLock()
	defer act.RUnlock()

	return act.StepIndex == len(act.Steps)
}

// StepTimeout once timed out will trigger an error and stop the automation.
// 2DO: Consider using context rather than a timeout.  Go programmers love context.
func (act *Action) StepTimeout() bool {
	act.RLock()
	defer act.RUnlock()

	b := false
	s := act.Steps[act.StepIndex]
	if s.Timeout > 0 {
		b = time.Now().After(act.Start.Add(s.Timeout))
	}
	return b
}

// ToJSON encodes the current step.  It will be sent to the server as a request.
func (act *Action) ToJSON() []byte {
	act.RLock()
	defer act.RUnlock()

	if act.Start == nil {
		t := time.Now()
		act.Start = &t
	}
	s := act.Steps[act.StepIndex]

	j, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return j
}

// Run sends the current action to websocket code that will create a request.
// Then the action will wait until all steps and expected events are completed.
func (act *Action) Run() error {
	ActionChan <- act
	return act.wait(ActionChan, Cache, StepChan)
}

// Log writes the current state of the action to the log.
func (act *Action) log() {
	log.Printf("Act %+v\n", act)
	for _, step := range act.Steps {
		log.Printf("Step %d Params %+v", step.ID, step.Params)
		log.Printf("Step %d Return %+v", step.ID, step.Reply)
	}
}
