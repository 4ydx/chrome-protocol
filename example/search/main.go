package main

import (
	"github.com/4ydx/chrome-protocol"
	"github.com/4ydx/chrome-protocol/actions/dom"
	"github.com/4ydx/chrome-protocol/actions/enable"
	"github.com/4ydx/chrome-protocol/actions/page"
	"log"
	"time"
)

func main() {
	id := cdp.Start()

	// Enable dom and page communication with chrome
	ea.EnableDom(id, time.Second*2).Run()
	ea.EnablePage(id, time.Second*2).Run()

	// Navigate
	a0 := pa.Navigate(id, "https://google.com", time.Second*5)
	a0.Run()

	// Focus
	if err := da.Focus(id, "#lst-ib", time.Second*5); err != nil {
		panic(err)
	}

	log.Print("\n-- All completed --\n")
	a0.Log()

	cdp.Stop()
}
