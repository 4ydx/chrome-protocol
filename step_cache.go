package main

import (
	"github.com/chromedp/cdproto"
	"log"
	"sync"
)

type StepCache struct {
	*sync.RWMutex
	*Step
}

func NewStepCache() *StepCache {
	cache := &StepCache{
		RWMutex: &sync.RWMutex{},
	}
	return cache
}

func (ac *StepCache) Set(s *Step) {
	ac.Lock()
	log.Printf("Set step %+v\n", s)
	ac.Step = s
	ac.Unlock()
}

func (ac *StepCache) GetId() int64 {
	ac.Lock()
	id := ac.Step.Id
	ac.Unlock()
	return id
}

func (ac *StepCache) SetResult(m cdproto.Message) int64 {
	ac.Lock()
	defer ac.Unlock()

	err := ac.Step.Returns.UnmarshalJSON(m.Result)
	if err != nil {
		log.Fatal("Unmarshal error:", err)
	}
	log.Printf(".RES: %+v\n", ac.Step)
	log.Printf("    : %+v\n", ac.Step.Params)
	log.Printf("    : %+v\n", ac.Step.Returns)
	id := ac.Step.ActionId

	return id
}
