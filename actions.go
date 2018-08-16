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
	// The id of the parent that contains this step.
	ActionId int64 `json:"-"`

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
	Id        int64
	Events    []Event
	Steps     []Step
	StepIndex int
	Start     *time.Time
}

type Actions []*Action

func (acts *Actions) Add(action *Action) {
	action.Id = int64(len(*acts))
	for i := 0; i < len(action.Steps); i++ {
		action.Steps[i].ActionId = action.Id
	}
	*acts = append(*acts, action)
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

func (act *Action) Wait(actions chan<- *Action, ec *EventCache, stepComplete <-chan int64) {
	for {
		select {
		case <-time.After(Wait):
			if !act.IsComplete() && act.StepTimeout() {
				log.Fatalf("Action %+v step timeout %+v\n", act, act.Step())
			}
			if act.IsComplete() && ec.EventsComplete() {
				log.Printf("Action %d completed.", act.Id)
				return
			}
			log.Printf("Action %d waiting...", act.Id)
		case id := <-stepComplete:
			act.Lock()
			if id != act.Id {
				log.Fatalf("Mismatched id %d != %d\n", id, act.Id)
			}
			log.Printf("Step %d complete with %+v", act.Steps[act.StepIndex].Id, act.Steps[act.StepIndex].Returns)
			act.StepIndex++
			act.Unlock()

			if !act.IsComplete() {
				actions <- act
			}
			if act.IsComplete() && ec.EventsComplete() {
				log.Printf("Action %d completed.", act.Id)
				return
			}
			log.Printf("Action %d waiting...", act.Id)
		}
	}
}
