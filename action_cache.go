package cdp

import (
	"github.com/4ydx/cdproto"
	"log"
)

type ActionCache struct {
	*Action
}

func NewActionCache() *ActionCache {
	return &ActionCache{}
}

func (ac *ActionCache) Set(a *Action) {
	a.RLock()
	defer a.RUnlock()

	log.Printf("Set action %+v\n", a)
	ac.Action = a
}

func (ac *ActionCache) HasStepId(id int64) bool {
	if ac.Action == nil {
		panic("Nil pointer")
	}
	ac.Action.RLock()
	defer ac.Action.RUnlock()

	if ac.Action.StepIndex == len(ac.Action.Steps) {
		return false
	}
	return ac.Action.Steps[ac.Action.StepIndex].Id == id
}

func (ac *ActionCache) SetResult(m cdproto.Message) {
	if ac.Action == nil {
		panic("Nil pointer")
	}
	ac.Action.Lock()
	defer ac.Action.Unlock()

	s := ac.Action.Steps[ac.Action.StepIndex]
	err := s.Returns.UnmarshalJSON(m.Result)
	if err != nil {
		log.Fatal("Unmarshal error:", err)
	}
	ac.Action.StepIndex++

	log.Printf(".STP COMPLETE: %+v\n", s)
	log.Printf("             : %+v\n", s.Params)
	log.Printf("             : %+v\n", s.Returns)
}
