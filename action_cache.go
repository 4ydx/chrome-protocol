package cdp

import (
	"github.com/4ydx/cdproto"
	"log"
)

// ActionCache stores the current action for safe use across routines.
type ActionCache struct {
	a *Action
}

func NewActionCache() *ActionCache {
	return &ActionCache{}
}

func (ac *ActionCache) Set(a *Action) {
	a.RLock()
	defer a.RUnlock()

	log.Printf("Set action %+v\n", a)
	ac.a = a
}

func (ac *ActionCache) HasStepId(id int64) bool {
	if ac.a == nil {
		log.Fatal("Nil pointer")
	}
	ac.a.RLock()
	defer ac.a.RUnlock()

	if ac.a.StepIndex == len(ac.a.Steps) {
		return false
	}
	return ac.a.Steps[ac.a.StepIndex].Id == id
}

func (ac *ActionCache) HasEvent(name cdproto.MethodType) bool {
	ac.a.Lock()
	defer ac.a.Unlock()

	_, ok := ac.a.Events[string(name)]
	return ok
}

func (ac *ActionCache) SetEventResult(name cdproto.MethodType, m cdproto.Message) {
	ac.a.Lock()
	defer ac.a.Unlock()

	if e, ok := ac.a.Events[string(name)]; ok {
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
		ac.a.Events[string(name)] = e

		log.Printf(".EVT: %s %+v\n", name, m)
		log.Printf("    : %+v\n", e)
		log.Printf("    : %+v\n", e.Value)
	}
}

func (ac *ActionCache) SetResult(m cdproto.Message) {
	if ac.a == nil {
		log.Fatal("Nil pointer")
	}
	ac.a.Lock()
	defer ac.a.Unlock()

	s := ac.a.Steps[ac.a.StepIndex]
	err := s.Returns.UnmarshalJSON(m.Result)
	if err != nil {
		log.Fatal("Unmarshal error:", err)
	}
	ac.a.StepIndex++

	log.Printf(".STP COMPLETE: %+v\n", s)
	log.Printf("             : %+v\n", s.Params)
	log.Printf("             : %+v\n", s.Returns)
}

func (ac *ActionCache) EventsComplete() bool {
	ac.a.RLock()
	defer ac.a.RUnlock()

	complete := true
	for _, e := range ac.a.Events {
		if e.IsRequired && !e.IsFound {
			complete = false
		}
	}
	return complete
}
