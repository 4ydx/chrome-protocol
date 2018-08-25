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

// IsStepComplete indicates that all steps are complete.
func (act *Action) IsStepComplete() bool {
	act.RLock()
	defer act.RUnlock()

	return act.StepIndex == len(act.Steps)
}

// IsComplete indicates that all steps and events are complete.
func (act *Action) IsComplete() bool {
	act.RLock()
	defer act.RUnlock()

	complete := true
	for _, e := range act.Events {
		if e.IsRequired && !e.IsFound {
			complete = false
		}
	}
	return act.StepIndex == len(act.Steps) && complete
}

// StepTimeout once timed out will trigger an error and stop the automation.
func (act *Action) StepTimeout() <-chan time.Time {
	act.RLock()
	defer act.RUnlock()

	return time.After(act.Steps[act.StepIndex].Timeout)
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
		case <-CacheCompleteChan:
			return nil
		case <-StepChan:
			ActionChan <- act.ToJSON()
			stepTimeout = act.StepTimeout()
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

// HasStepID determines if an id matches the current action's step's unique id.
func (act *Action) HasStepID(id int64) bool {
	act.RLock()
	defer act.RUnlock()

	if act.StepIndex == len(act.Steps) {
		return false
	}
	return act.Steps[act.StepIndex].ID == id
}

// HasEvent returns true when the action has an event with the given MethodType.
func (act *Action) HasEvent(name string) bool {
	act.RLock()
	defer act.RUnlock()

	_, ok := act.Events[name]
	return ok
}

// GetStepMethod returns the method of the step that is currently active or the very last method.
func (act *Action) GetStepMethod() string {
	act.RLock()
	defer act.RUnlock()

	if act.StepIndex == len(act.Steps) {
		return act.Steps[act.StepIndex-1].Method
	}
	return act.Steps[act.StepIndex].Method
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
					log.Printf("No matching frameID")
					return nil
				}
			} else {
				if ok := e.Value.MatchFrameID(frameID, m.Result); !ok {
					log.Printf("No matching frameID")
					return nil
				}
			}
		}
		e.IsFound = true
		act.Events[string(name)] = e

		log.Printf(".EVT: %s %+v\n", name, m)
		log.Printf("    : %+v\n", e)
		log.Printf("    : %+v\n", e.Value)
	}
	return nil
}

// SetResult applies the message returns to the current step and advances the step.
func (act *Action) SetResult(m Message) error {
	act.Lock()
	defer act.Unlock()

	s := act.Steps[act.StepIndex]
	frameID := act.Frame.GetFrameID()
	if frameID == "" {
		err := s.Reply.UnmarshalJSON(m.Result)
		if err != nil {
			log.Fatalf("Unmarshal error: %s", err)
		}
		act.Frame.SetFrameID(s.Reply.GetFrameID())
	} else {
		if ok := s.Reply.MatchFrameID(frameID, m.Result); !ok {
			log.Printf("No matching frameID")
			return nil
		}
	}
	act.StepIndex++

	log.Printf(".STP COMPLETE: %+v\n", s)
	log.Printf("             : %+v\n", s.Params)
	log.Printf("             : %+v\n", s.Reply)

	return nil
}
