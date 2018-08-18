package cdp

import (
	"errors"
	"fmt"
	"github.com/4ydx/cdp/protocol/dom"
	"sync"
)

type Page struct {
	*sync.RWMutex
	DOM      *dom.GetFlattenedDocumentReply
	FrameId  string
	LoaderId string
	ID       ID
}

func (p *Page) CheckFrameId(pi *ProtocolIds) error {
	p.Lock()
	defer p.Unlock()

	if p.FrameId == "" {
		p.FrameId = pi.FID
		p.FrameId = pi.FID
	} else if p.FrameId != pi.FID {
		return errors.New(fmt.Sprintf("FrameID mismatch %s != %s.", p.FrameId, pi.FID))
	}
	return nil
}
