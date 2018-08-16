package cdp

import (
	"encoding/json"
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
	Events    []Event
	Steps     []Step
	StepIndex int
	Start     *time.Time
}

func NewAction(events []Event, steps []Step) *Action {
	return &Action{
		RWMutex: &sync.RWMutex{},
		Events:  events,
		Steps:   steps,
	}
}

func (act *Action) IsComplete() bool {
	act.RLock()
	defer act.RUnlock()

	return act.StepIndex == len(act.Steps)
}

func (act *Action) StepTimeout() bool {
	b := false
	act.RLock()
	s := act.Steps[act.StepIndex]
	if s.Timeout > 0 {
		b = time.Now().After(act.Start.Add(s.Timeout))
	}
	act.RUnlock()
	return b
}

func (act *Action) Step() Step {
	act.Lock()
	if act.Start == nil {
		t := time.Now()
		act.Start = &t
	}
	s := act.Steps[act.StepIndex]
	act.Unlock()
	return s
}

func (act *Action) ToJSON() []byte {
	j, err := json.Marshal(act.Step())
	if err != nil {
		log.Fatal(err)
	}
	return j
}

func (act *Action) Wait(actions chan<- *Action, ec *EventCache, stepComplete <-chan bool) {
	for {
		select {
		case <-time.After(Wait):
			if !act.IsComplete() && act.StepTimeout() {
				log.Fatalf("Action %+v step timeout %+v\n", act, act.Step())
			}
			if act.IsComplete() && ec.EventsComplete() {
				log.Print("Action completed.")
				return
			}
			log.Print("Action waiting...")
		case <-stepComplete:
			act.Lock()
			log.Printf("Step %d complete with %+v", act.Steps[act.StepIndex].Id, act.Steps[act.StepIndex].Returns)
			act.StepIndex++
			act.Unlock()

			if !act.IsComplete() {
				actions <- act
			}
			if act.IsComplete() && ec.EventsComplete() {
				log.Printf("Action completed.")
				return
			}
			log.Printf("Action waiting...")
		}
	}
}

func (act *Action) Run(ec *EventCache, actionChan chan<- *Action, stepComplete <-chan bool) {
	ec.Load(act.Events)
	actionChan <- act
	act.Wait(actionChan, ec, stepComplete)
	ec.Log()
}

func (act *Action) Log() {
	log.Printf("Act %+v\n", act)
	for _, step := range act.Steps {
		log.Printf("Step %d Params %+v", step.Id, step.Params)
		log.Printf("Step %d Return %+v", step.Id, step.Returns)
	}
}
