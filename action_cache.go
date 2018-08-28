package cdp

import (
	"log"
	"sync"
	"time"
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

// HasCommandID determines if an id matches the current action's command's unique id.
func (ac *ActionCache) HasCommandID(id int64) bool {
	ac.RLock()
	defer ac.RUnlock()

	if ac.a == nil {
		return false
	}
	return ac.a.HasCommandID(id)
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

// GetCommandMethod returns the method of the command that is currently active.
func (ac *ActionCache) GetCommandMethod() string {
	ac.RLock()
	defer ac.RUnlock()

	if ac.a == nil {
		return ""
	}
	return ac.a.GetCommandMethod()
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

// SetResult applies the message returns to the current command and advances the command.
func (ac *ActionCache) SetResult(m Message) error {
	ac.Lock()
	defer ac.Unlock()

	if ac.a == nil {
		return nil
	}
	return ac.a.SetResult(m)
}

// IsComplete indicates whether or not all events and commands are completed.
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

// IsCommandComplete indicates whether or not the command portion of an action is completed.
func (ac *ActionCache) IsCommandComplete() bool {
	ac.RLock()
	defer ac.RUnlock()

	if ac.a == nil {
		return true
	}
	if ac.a.IsCommandComplete() {
		return true
	}
	return false
}

// CommandTimeout returns the timeout channel for the current command.
func (ac *ActionCache) CommandTimeout() <-chan time.Time {
	ac.RLock()
	defer ac.RUnlock()

	if ac.a == nil {
		return time.After(0)
	}
	return ac.a.CommandTimeout()
}

// ToJSON returns the json representation of the current command.
func (ac *ActionCache) ToJSON() []byte {
	ac.RLock()
	defer ac.RUnlock()

	if ac.a == nil {
		return []byte("")
	}
	return ac.a.ToJSON()
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
