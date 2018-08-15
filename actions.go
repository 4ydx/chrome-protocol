package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

var Wait = time.Millisecond * 50

type Step struct {
	ActionId int64 `json:"-"`

	Id      int64            `json:"id"`
	Method  string           `json:"method"`
	Params  json.Marshaler   `json:"params"`
	Returns json.Unmarshaler `json:"-"`

	Timeout time.Duration `json:"-"`
}

type Action struct {
	*sync.RWMutex
	Id        int64
	Steps     []*Step
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

func NewAction(steps []*Step) *Action {
	return &Action{
		RWMutex: &sync.RWMutex{},
		Steps:   steps,
	}
}

func (act *Action) IsComplete() bool {
	act.RLock()
	b := act.StepIndex == len(act.Steps)
	act.RUnlock()
	return b
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

func (act *Action) Step() *Step {
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

func (act *Action) Wait(stepComplete <-chan int64) {
	for {
		select {
		case <-time.After(Wait):
			if act.IsComplete() {
				log.Printf("Action %d completed.", act.Id)
				return
			}
			if act.StepTimeout() {
				log.Fatalf("Action %+v step timeout %+v\n", act, act.Step())
			}
			log.Print("waiting...")
		case id := <-stepComplete:
			act.Lock()
			if id != act.Id {
				log.Fatalf("Mismatched id %d != %d\n", id, act.Id)
			}
			log.Printf("Step %d complete", act.Steps[act.StepIndex].Id)
			act.StepIndex++
			act.Unlock()
		}
	}
}
