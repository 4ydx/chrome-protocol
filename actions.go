package cdp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var Wait = time.Millisecond * 50

type Event struct {
	Name       string
	Value      json.Unmarshaler
	IsRequired bool
	IsFound    bool
	OnEvent    func(Event) // Callback for accessing the event
}

type Step struct {
	// Values required to make a chrome devtools protocol request
	Id     int64          `json:"id"`
	Method string         `json:"method"`
	Params json.Marshaler `json:"params"`

	Returns         json.Unmarshaler `json:"-"` // The struct that will be filled when a matching step Id is found in a reply over the chrome websocket.
	Timeout         time.Duration    `json:"-"` // How long until the current step experiences a timeout, which will halt the entire process.
	PreviousReturns func()           `json:"-"` // Method defined by the user to take the previous step's Returns and apply them to the current step's Params.
}

type Action struct {
	*sync.RWMutex
	Steps     []Step
	StepIndex int
	Events    map[string]Event
	Start     *time.Time
}

func NewAction(events []Event, steps []Step) *Action {
	act := &Action{
		RWMutex: &sync.RWMutex{},
		Events:  make(map[string]Event),
		Steps:   steps,
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
				return errors.New(fmt.Sprintf("Action %s step timeout\n", act.ToJSON()))
			}
			if act.IsComplete() && ac.EventsComplete() {
				log.Print("Action completed.")
				return nil
			}
			log.Print("Action waiting...")
		case <-stepChan:
			if !act.IsComplete() {
				if act.StepTimeout() {
					return errors.New(fmt.Sprintf("Action %s step timeout\n", act.ToJSON()))
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

func (act *Action) Log() {
	log.Printf("Act %+v\n", act)
	for _, step := range act.Steps {
		log.Printf("Step %d Params %+v", step.Id, step.Params)
		log.Printf("Step %d Return %+v", step.Id, step.Returns)
	}
}
