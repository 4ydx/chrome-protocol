package main

import (
	"github.com/chromedp/cdproto"
	"log"
	"sync"
)

type EventCache struct {
	*sync.RWMutex
	Events map[string]*Event
}

func NewEventCache() *EventCache {
	cache := &EventCache{
		RWMutex: &sync.RWMutex{},
		Events:  make(map[string]*Event),
	}
	return cache
}

func (ac *EventCache) Add(s *Event) {
	ac.Lock()
	ac.Events[s.Name] = s
	ac.Unlock()
}

func (ac *EventCache) HasEvent(name string) (*Event, bool) {
	ac.Lock()
	e, ok := ac.Events[name]
	ac.Unlock()
	return e, ok
}

func (ac *EventCache) SetResult(name string, m cdproto.Message) {
	ac.Lock()
	defer ac.Unlock()

	if e, ok := ac.Events[name]; ok {
		err := e.Returns.UnmarshalJSON(m.Result)
		if err != nil {
			log.Fatal("Unmarshal error:", err)
		}
		log.Printf(".RES: %+v\n", e)
		log.Printf("    : %+v\n", e.Returns)
	}
}
