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
	ac.RLock()
	defer ac.RUnlock()

	if ac.Action.StepIndex == len(ac.Action.Steps) {
		return false
	}
	return ac.Action.Steps[ac.Action.StepIndex].Id == id
}

func (ac *ActionCache) HasEvent(name cdproto.MethodType) bool {
	ac.Lock()
	defer ac.Unlock()

	_, ok := ac.Events[string(name)]
	return ok
}

func (ac *ActionCache) SetEventResult(name cdproto.MethodType, m cdproto.Message) {
	ac.Lock()
	defer ac.Unlock()

	if e, ok := ac.Events[string(name)]; ok {
		err := e.Value.UnmarshalJSON(m.Params)
		if err != nil {
			log.Printf("Unmarshal error: %s; for %+v; from %+v", err.Error(), e.Value, m)
			err = e.Value.UnmarshalJSON(m.Result)
			if err != nil {
				log.Printf("Unmarshal error: %s; for %+v; from %+v", err.Error(), e.Value, m)
				return
			}
		}
		e.IsFound = true
		ac.Events[string(name)] = e

		log.Printf(".SET: %s %+v\n", name, m)
		log.Printf("    : %+v\n", e)
		log.Printf("    : %+v\n", e.Value)
	}
}

func (ac *ActionCache) SetResult(m cdproto.Message) {
	if ac.Action == nil {
		panic("Nil pointer")
	}
	ac.Lock()
	defer ac.Unlock()

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

func (ac *ActionCache) EventsComplete() bool {
	ac.RLock()
	defer ac.RUnlock()

	complete := true
	for _, e := range ac.Events {
		if e.IsRequired && !e.IsFound {
			complete = false
		}
	}
	return complete
}
