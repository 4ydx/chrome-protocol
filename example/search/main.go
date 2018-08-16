package main

import (
	"github.com/4ydx/chrome-protocol"
	"github.com/4ydx/chrome-protocol/actions/dom"
	"github.com/4ydx/chrome-protocol/actions/enable"
	"github.com/4ydx/chrome-protocol/actions/page"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	f, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(f)

	c := cdp.GetWebsocket()
	defer c.Close()

	shutdown := make(chan struct{})
	actions := make(chan *cdp.Action)
	stepComplete := make(chan bool)
	allComplete := make(chan struct{})

	stepCache := cdp.NewStepCache()
	eventCache := cdp.NewEventCache()

	go cdp.Write(c, actions, stepCache, shutdown, allComplete)
	go cdp.Read(c, stepComplete, stepCache, eventCache, shutdown)

	id := &cdp.ID{
		RWMutex: &sync.RWMutex{},
		Value:   11111,
	}

	var acts cdp.Actions

	// Enable communication with chrome
	acts.Add(ea.EnableDom(id, time.Second*2))
	acts.Add(ea.EnablePage(id, time.Second*2))

	// Navigate
	acts.Add(pa.Navigate(id, "https://google.com", time.Second*5))

	// Find the element
	acts.Add(da.Find(id, "lst-ib", time.Second*5))

	for i := 0; i < len(acts); i++ {
		eventCache.Load(acts[i].Events)
		actions <- acts[i]
		acts[i].Wait(actions, eventCache, stepComplete)
		eventCache.Log()
	}
	log.Print("\n-- All completed --\n")
	for _, act := range acts {
		log.Printf("Act %+v\n", act)
		for i, step := range act.Steps {
			log.Printf("Step %d Params %+v", i, step.Params)
			log.Printf("Step %d Return %+v", i, step.Returns)
		}
	}
	allComplete <- struct{}{}
	<-shutdown
}
