package main

import (
	"log"
	"sync"
)

type StepId int64

type StepCache struct {
	*sync.RWMutex
	Map map[StepId]*Step
}

func NewStepCache() *StepCache {
	cache := &StepCache{
		RWMutex: &sync.RWMutex{},
		Map:     make(map[StepId]*Step),
	}
	return cache
}

func (ac *StepCache) Add(a *Step) {
	ac.Lock()
	log.Printf("Add action %+v\n", a)
	ac.Map[StepId(a.Id)] = a
	ac.Unlock()
}

func (ac *StepCache) Get(id int64) (*Step, bool) {
	ac.Lock()
	a, ok := ac.Map[StepId(id)]
	ac.Unlock()
	return a, ok
}
