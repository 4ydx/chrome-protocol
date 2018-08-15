package main

import (
	"github.com/chromedp/cdproto"
	"log"
	"sync"
)

type EventCache struct {
	*sync.RWMutex
	Events map[cdproto.MethodType]*Event
}

func NewEventCache() *EventCache {
	cache := &EventCache{
		RWMutex: &sync.RWMutex{},
		Events:  make(map[cdproto.MethodType]*Event),
	}
	return cache
}

func (ec *EventCache) Log() {
	ec.RLock()
	defer ec.RUnlock()

	for _, event := range ec.Events {
		log.Printf("Event %+v\n", event)
		log.Printf("Event Value %+v", event.Value)
	}
}

func (ec *EventCache) Load(events []Event) {
	ec.Lock()
	defer ec.Unlock()

	ec.Events = make(map[cdproto.MethodType]*Event)
	for _, e := range events {
		ec.Events[cdproto.MethodType(e.Name)] = &e
	}
}

func (ec *EventCache) Add(s *Event) {
	ec.Lock()
	ec.Events[cdproto.MethodType(s.Name)] = s
	ec.Unlock()
}

func (ec *EventCache) HasEvent(name cdproto.MethodType) (*Event, bool) {
	ec.Lock()
	e, ok := ec.Events[name]
	ec.Unlock()
	return e, ok
}

func (ec *EventCache) EventsComplete() bool {
	ec.RLock()
	defer ec.RUnlock()

	complete := true
	for _, e := range ec.Events {
		if e.IsRequired && !e.IsFound {
			complete = false
		}
	}
	return complete
}

func (ec *EventCache) SetResult(name cdproto.MethodType, m cdproto.Message) {
	ec.Lock()
	defer ec.Unlock()

	if e, ok := ec.Events[name]; ok {
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
		ec.Events[name] = e

		log.Printf(".RES: %+v\n", e)
		log.Printf("    : %+v\n", e.Value)
	}
}
