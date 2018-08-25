package cdp

import (
	"log"
	"sync"
)

// ActionCache stores the current action for safe use across routines.
type ActionCache struct {
	*sync.RWMutex
	a *Action
}

// Set puts the given action into the action cache.
// Subsequent incoming request processing will use this cached action.
func (ac *ActionCache) Set(a *Action) {
	ac.Lock()
	defer ac.Unlock()

	log.Printf(".SET action %+v\n", a)
	ac.a = a
}

// HasStepID determines if an id matches the current action's step's unique id.
func (ac *ActionCache) HasStepID(id int64) bool {
	ac.RLock()
	defer ac.RUnlock()

	if ac.a == nil {
		return false
	}
	return ac.a.HasStepID(id)
}

// HasEvent returns true when the action has an event with the given MethodType.
func (ac *ActionCache) HasEvent(name string) bool {
	ac.RLock()
	defer ac.RUnlock()

	if ac.a == nil {
		return false
	}
	return ac.a.HasEvent(name)
}

// GetStepMethod returns the method of the step that is currently active.
func (ac *ActionCache) GetStepMethod() string {
	ac.RLock()
	defer ac.RUnlock()

	if ac.a == nil {
		return ""
	}
	return ac.a.GetStepMethod()
}

// GetFrameID returns the frameID of the current frame.
func (ac *ActionCache) GetFrameID() string {
	ac.Lock()
	defer ac.Unlock()

	if ac.a == nil {
		return ""
	}
	return ac.a.GetFrameID()
}

// SetEvent takes the given message and sets an event's params or results's.
func (ac *ActionCache) SetEvent(name string, m Message) error {
	ac.Lock()
	defer ac.Unlock()

	if ac.a == nil {
		return nil
	}
	return ac.a.SetEvent(name, m)
}

// SetResult applies the message returns to the current step and advances the step.
func (ac *ActionCache) SetResult(m Message) error {
	ac.Lock()
	defer ac.Unlock()

	if ac.a == nil {
		return nil
	}
	return ac.a.SetResult(m)
}

// IsComplete indicates whether or not all events and steps are completed.
func (ac *ActionCache) IsComplete() bool {
	ac.RLock()
	defer ac.RUnlock()

	if ac.a == nil {
		return true
	}
	if ac.a.IsComplete() {
		return true
	}
	return false
}

// IsStepComplete indicates whether or not the step portion of an action is completed.
func (ac *ActionCache) IsStepComplete() bool {
	ac.RLock()
	defer ac.RUnlock()

	if ac.a == nil {
		return true
	}
	if ac.a.IsStepComplete() {
		return true
	}
	return false
}

// Clear the cached action.
func (ac *ActionCache) Clear() {
	ac.Lock()
	defer ac.Unlock()

	if ac.a == nil {
		return
	}
	ac.a = nil
}
