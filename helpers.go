package cdp

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Message is the chrome DevTools Protocol message sent/read over the websocket connection.
type Message struct {
	ID     int64           `json:"id,omitempty"`     // Unique message identifier.
	Method string          `json:"method,omitempty"` // Event or command type.
	Params json.RawMessage `json:"params,omitempty"` // Event or command parameters.
	Result json.RawMessage `json:"result,omitempty"` // Command return values.
	Error  *Error          `json:"error,omitempty"`  // Error message.
}

// Error error type that is apart of the Message struct.
type Error struct {
	Code    int64  `json:"code"`    // Error code.
	Message string `json:"message"` // Error message.
}

// Error satisfies the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

// RequestID stores the last value used for chrome devtool protocal requests being sent to the server.
type RequestID struct {
	*sync.RWMutex
	Value int64
}

// GetNext is a convenience method for creating the unique ids required when performing chrome devtool protocol requests.
func (id *RequestID) GetNext() int64 {
	id.Lock()
	id.Value++
	v := id.Value
	id.Unlock()
	return v
}
