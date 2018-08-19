package cdp

import (
	"fmt"
	"github.com/4ydx/cdp/protocol/dom"
	"sync"
)

// Page stores the current FrameID.
type Page struct {
	*sync.RWMutex
	DOM      *dom.GetFlattenedDocumentReply
	FrameID  string
	LoaderID string
	ID       ID
}

// CheckFrameID attempts to validate the FrameID.
// This will likely change.
func (p *Page) CheckFrameID(pi *ProtocolIds) error {
	p.Lock()
	defer p.Unlock()

	if p.FrameID == "" {
		p.FrameID = pi.FID
		p.FrameID = pi.FID
	} else if p.FrameID != pi.FID {
		return fmt.Errorf("frameid mismatch %s != %s", p.FrameID, pi.FID)
	}
	return nil
}

func (p *Page) GetFrameID() string {
	p.RLock()
	defer p.RUnlock()

	return p.FrameID
}

func (p *Page) SetFrameID(frameID string) {
	p.Lock()
	defer p.Unlock()

	p.FrameID = frameID
}
