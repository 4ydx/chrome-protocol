package main

import (
	"github.com/chromedp/cdproto/dom"
)

type FindNodes struct {
	Step1       *dom.PerformSearchParams
	Step1Return *dom.PerformSearchReturns
	Step2       *dom.GetSearchResultsParams
	Step2Return *dom.GetSearchResultsReturns
	Complete    bool
}

func (fn *FindNodes) All() {
	// continuously run the find all method until all nodes are found
}
