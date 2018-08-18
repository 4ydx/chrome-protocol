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

	log.Printf("Set action %+v\n", a)
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

// SetEvent takes the given message and sets an event's params or results's.
func (ac *ActionCache) SetEvent(name string, m Message, pi *ProtocolIds) error {
	ac.a.Lock()
	defer ac.a.Unlock()

	if err := ac.a.Page.CheckFrameId(pi); err != nil {
		// When the frame ID differs, it indicates that this Event is not intended for this Action.
		// In other words a different Action needs to consume this event.
		return nil
	}
	if e, ok := ac.a.Events[name]; ok {
		err := e.Value.UnmarshalJSON(m.Params)
		if err != nil {
			log.Printf("Unmarshal error: %s; for %+v; from %+v", err.Error(), e.Value, m)
			err = e.Value.UnmarshalJSON(m.Result)
			if err != nil {
				log.Printf("Unmarshal error: %s; for %+v; from %+v", err.Error(), e.Value, m)
				return err
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
func (ac *ActionCache) SetResult(m Message, pi *ProtocolIds) error {
	if ac.a == nil {
		log.Fatal("Nil pointer")
	}
	ac.a.Lock()
	defer ac.a.Unlock()

	if err := ac.a.Page.CheckFrameId(pi); err != nil {
		ac.a.log()
		return err
	}
	s := ac.a.Steps[ac.a.StepIndex]
	err := s.Reply.UnmarshalJSON(m.Result)
	if err != nil {
		log.Fatalf("Unmarshal error: %s", err)
	}

	// TODO: Check to see that the frameID is correct. if not, clear the object of the marshaling changes and do not increment the step.
	// NEED: a CheckFrameID method.  The Returns needs to have Unmarshalling as well as a new CheckFrameID method...

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
