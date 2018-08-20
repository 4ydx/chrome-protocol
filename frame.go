package cdp

import (
	"github.com/4ydx/cdp/protocol/dom"
	"sync"
)

// Frame stores the current FrameID.
type Frame struct {
	*sync.RWMutex
	DOM       *dom.GetFlattenedDocumentReply
	FrameID   string
	LoaderID  string
	RequestID RequestID
}

func (p *Frame) GetFrameID() string {
	p.RLock()
	defer p.RUnlock()

	return p.FrameID
}

func (p *Frame) SetFrameID(frameID string) {
	p.Lock()
	defer p.Unlock()

	p.FrameID = frameID
}