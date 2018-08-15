package main

import (
	lg "github.com/chromedp/cdproto/log"
	"github.com/chromedp/cdproto/page"
	"log"
	"os"
	"sync"
	"time"
	//"github.com/chromedp/cdproto/browser"
	//"github.com/chromedp/cdproto/css"
	//"github.com/chromedp/cdproto/dom"
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
	cache := NewStepCache()

	go Write(c, actions, cache, shutdown, allComplete)
	go Read(c, stepComplete, cache, shutdown)

	id := &ID{
		RWMutex: &sync.RWMutex{},
		Value:   11111,
	}

	var process Actions
	process.Add(
		NewAction([]*Step{&Step{Id: id.GetNext(), Method: lg.CommandEnable, Params: &lg.EnableParams{}, Returns: &lg.EnableReturns{}, Timeout: time.Second * 3}}),
	)
	process.Add(
		NewAction([]*Step{&Step{Id: id.GetNext(), Method: page.CommandEnable, Params: &page.EnableParams{}, Returns: &page.EnableReturns{}, Timeout: time.Second * 3}}),
	)
	process.Add(
		NewAction([]*Step{
			&Step{
				Id:      id.GetNext(),
				Method:  page.CommandNavigate,
				Params:  &page.NavigateParams{URL: "https://google.com"},
				Returns: &page.NavigateReturns{},
				Timeout: time.Second * 10,
			},
		}),
	)
	//Action{Id: id.GetNext(), Method: runtime.CommandEnable, Wait: time.Second * 0},
	//Action{Id: id.GetNext(), Method: inspector.CommandEnable, Wait: time.Second * 0},
	//Action{Id: id.GetNext(), Method: page.CommandEnable, Wait: time.Second * 0},
	//Action{Id: id.GetNext(), Method: dom.CommandEnable, Wait: time.Second * 0},
	//Action{Id: id.GetNext(), Method: css.CommandEnable, Wait: time.Second * 0},
	//Action{Id: id.GetNext(), Method: network.CommandEnable, Wait: time.Second * 0},

	//Action{Id: id.GetNext(), Method: page.CommandReload, Wait: time.Second * 5},
	//Action{Id: id.GetNext(), Method: page.CommandNavigate, Params: page.NavigateParams{URL: "https://google.com"}, Wait: time.Second * 20},

	//Action{Id: id.GetNext(), Method: dom.CommandPerformSearch, Params: dom.PerformSearchParams{Query: "#login_form_2 input[name='Email']"}, Wait: time.Second * 10},

	//Action{Id: id.GetNext(), Method: browser.CommandClose, Wait: time.Second * 5},
	for i := 0; i < len(process); i++ {
		actions <- process[i]
		process[i].Wait(stepComplete)
	}
	log.Print("All completed.")
	allComplete <- struct{}{}
	<-shutdown
}
