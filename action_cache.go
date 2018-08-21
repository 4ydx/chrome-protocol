package cdp

import (
	"log"
)

// ActionCache stores the current action for safe use across routines.
type ActionCache struct {
	a *Action
}

// Set puts the given action into the action cache.
// Subsequent incoming request processing will use this cached action.
func (ac *ActionCache) Set(a *Action) {
	a.Lock()
	defer a.Unlock()

	log.Printf(".SET action %+v\n", a)
	ac.a = a
}

// HasStepID determines if an id matches the current action's step's unique id.
func (ac *ActionCache) HasStepID(id int64) bool {
	if ac.a == nil {
		log.Fatal("Nil pointer")
	}
	ac.a.RLock()
	defer ac.a.RUnlock()

	if ac.a.StepIndex == len(ac.a.Steps) {
		return false
	}
	return ac.a.Steps[ac.a.StepIndex].ID == id
}

// HasEvent returns true when the action has an event with the given MethodType.
func (ac *ActionCache) HasEvent(name string) bool {
	ac.a.Lock()
	defer ac.a.Unlock()

	_, ok := ac.a.Events[name]
	return ok
}

// GetStepMethod returns the method of the step that is currently active.
func (ac *ActionCache) GetStepMethod() string {
	ac.a.Lock()
	defer ac.a.Unlock()
	if ac.a.StepIndex == len(ac.a.Steps) {
		return ac.a.Steps[ac.a.StepIndex-1].Method
	}
	return ac.a.Steps[ac.a.StepIndex].Method
}

// GetFrameID returns the frameID of the current frame.
func (ac *ActionCache) GetFrameID() string {
	ac.a.Lock()
	defer ac.a.Unlock()
	return ac.a.Frame.FrameID
}

// SetEvent takes the given message and sets an event's params or results's.
func (ac *ActionCache) SetEvent(name string, m Message) error {
	ac.a.Lock()
	defer ac.a.Unlock()

	// Attempt to compare the incoming Event's frameID value with the existing value.
	frameID := ac.a.Frame.GetFrameID()
	if e, ok := ac.a.Events[name]; ok {
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
		ac.a.Events[string(name)] = e

		log.Printf(".EVT: %s %+v\n", name, m)
		log.Printf("    : %+v\n", e)
		log.Printf("    : %+v\n", e.Value)
	}
	return nil
}

// SetResult applies the message returns to the current step and advances the step.
func (ac *ActionCache) SetResult(m Message) error {
	if ac.a == nil {
		log.Fatal("Nil pointer")
	}
	ac.a.Lock()
	defer ac.a.Unlock()

	s := ac.a.Steps[ac.a.StepIndex]
	frameID := ac.a.Frame.GetFrameID()
	if frameID == "" {
		err := s.Reply.UnmarshalJSON(m.Result)
		if err != nil {
			log.Fatalf("Unmarshal error: %s", err)
		}
		ac.a.Frame.SetFrameID(s.Reply.GetFrameID())
	} else {
		if ok := s.Reply.MatchFrameID(frameID, m.Result); !ok {
			log.Printf("No matching frameID")
			return nil
		}
	}
	ac.a.StepIndex++

	log.Printf(".STP COMPLETE: %+v\n", s)
	log.Printf("             : %+v\n", s.Params)
	log.Printf("             : %+v\n", s.Reply)

	return nil
}

// EventsComplete indicates whether or not all required events have received a message from the server.
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
