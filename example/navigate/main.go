package main

import (
	"github.com/chromedp/cdproto"
	"github.com/chromedp/cdproto/dom"
	lg "github.com/chromedp/cdproto/log"
	"github.com/chromedp/cdproto/page"
	"log"
	"os"
	"sync"
	"time"
	//"github.com/chromedp/cdproto/browser"
	//"github.com/chromedp/cdproto/css"
	//"github.com/chromedp/cdproto/inspector"
	//"github.com/chromedp/cdproto/runtime"
	//"github.com/chromedp/cdproto/network"
	//"github.com/gorilla/websocket"
)

type ID struct {
	*sync.RWMutex
	Value int64
}

func (id *ID) GetNext() int64 {
	id.Lock()
	id.Value += 1
	v := id.Value
	id.Unlock()
	return v
}

func main() {
	f, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(f)

	c := GetWebsocket()
	defer c.Close()

	shutdown := make(chan struct{})
	actions := make(chan *Action)
	stepComplete := make(chan int64)
	allComplete := make(chan struct{})

	stepCache := NewStepCache()
	eventCache := NewEventCache()

	go Write(c, actions, stepCache, shutdown, allComplete)
	go Read(c, stepComplete, stepCache, eventCache, shutdown)

	id := &ID{
		RWMutex: &sync.RWMutex{},
		Value:   11111,
	}

	var acts Actions

	// Tell Chrome Devtools Protocol what data to send
	acts.Add(
		NewAction([]Event{}, []Step{
			Step{Id: id.GetNext(), Method: lg.CommandEnable, Params: &lg.EnableParams{}, Returns: &lg.EnableReturns{}, Timeout: time.Second * 3},
			Step{Id: id.GetNext(), Method: page.CommandEnable, Params: &page.EnableParams{}, Returns: &page.EnableReturns{}, Timeout: time.Second * 3},
			Step{Id: id.GetNext(), Method: dom.CommandEnable, Params: &dom.EnableParams{}, Returns: &dom.EnableReturns{}, Timeout: time.Second * 3},
		}),
	)

	// Navigate to a url, waiting for the page to stop loading
	acts.Add(
		NewAction(
			[]Event{
				Event{Name: cdproto.EventPageFrameStartedLoading, Value: &page.EventFrameStartedLoading{}, IsRequired: true},
				Event{Name: cdproto.EventPageFrameStoppedLoading, Value: &page.EventFrameStoppedLoading{}, IsRequired: true},
			},
			[]Step{
				Step{Id: id.GetNext(), Method: page.CommandNavigate, Params: &page.NavigateParams{URL: "https://google.com"}, Returns: &page.NavigateReturns{}, Timeout: time.Second * 10},
			}),
	)

	// TODO: Searching the DOM - will have to have a way to pass values between steps...

	//Action{Id: id.GetNext(), Method: runtime.CommandEnable, Wait: time.Second * 0},
	//Action{Id: id.GetNext(), Method: inspector.CommandEnable, Wait: time.Second * 0},
	//Action{Id: id.GetNext(), Method: css.CommandEnable, Wait: time.Second * 0},
	//Action{Id: id.GetNext(), Method: network.CommandEnable, Wait: time.Second * 0},

	//Action{Id: id.GetNext(), Method: page.CommandReload, Wait: time.Second * 5},
	//Action{Id: id.GetNext(), Method: page.CommandNavigate, Params: page.NavigateParams{URL: "https://google.com"}, Wait: time.Second * 20},

	//Action{Id: id.GetNext(), Method: dom.CommandPerformSearch, Params: dom.PerformSearchParams{Query: "#login_form_2 input[name='Email']"}, Wait: time.Second * 10},

	//Action{Id: id.GetNext(), Method: browser.CommandClose, Wait: time.Second * 5},

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
