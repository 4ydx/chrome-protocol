package cdp

import (
	"sync"
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
