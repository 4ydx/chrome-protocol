package cdp

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// Message chrome DevTools Protocol message sent/read over websocket
// connection.
type Message struct {
	ID     int64           `json:"id,omitempty"`     // Unique message identifier.
	Method string          `json:"method,omitempty"` // Event or command type.
	Params json.RawMessage `json:"params,omitempty"` // Event or command parameters.
	Result json.RawMessage `json:"result,omitempty"` // Command return values.
	Error  *Error          `json:"error,omitempty"`  // Error message.
}

// Error error type.
type Error struct {
	Code    int64  `json:"code"`    // Error code.
	Message string `json:"message"` // Error message.
}

// Error satisfies the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

// ProtocolIds contains the frame ID and loader ID from Params/Result being returned by the chrome devtools protocol.
type ProtocolIds struct {
	FID string `json:"frameId"`
	LID string `json:"loaderId"`
}

// UnmarshalIds extracts the FrameID and LoaderID from the message received from the server.
func UnmarshalIds(m Message) (*ProtocolIds, error) {
	pi := &ProtocolIds{}
	err := json.Unmarshal(m.Result, pi)
	if err != nil {
		log.Printf("Failure to unmarshal result %+v into ProtocolIds", m.Result)
		err = json.Unmarshal(m.Params, pi)
	}
	return pi, err
}

// ID stores the last value used for chrome devtool protocal requests being sent to the server.
type ID struct {
	*sync.RWMutex
	Value int64
}

// GetNext is a convenience method for creating the unique ids required when performing chrome devtool protocol requests.
func (id *ID) GetNext() int64 {
	id.Lock()
	id.Value++
	v := id.Value
	id.Unlock()
	return v
}
