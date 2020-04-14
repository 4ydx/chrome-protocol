package cdp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/4ydx/cdp/protocol/dom"
	"github.com/4ydx/cdp/protocol/lib"
	"github.com/gorilla/websocket"
)

// GetWebsocket returns a websocket connection to the running browser.
func GetWebsocket(lg *log.Logger, port int) *websocket.Conn {
	r, err := http.Get(fmt.Sprintf("http://localhost:%d/json", port))
	if err != nil {
		lg.Fatal(err)
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		lg.Fatal(err)
	}
	var inter interface{}
	err = json.Unmarshal(b, &inter)
	if err != nil {
		lg.Fatal(err)
	}

	ws := ""
	arr := inter.([]interface{})
	for i := 0; i < len(arr); i++ {
		entry := arr[i].(map[string]interface{})
		for k, v := range entry {
			if k == "webSocketDebuggerUrl" {
				ws = v.(string)
			}
		}
	}
	if ws == "" {
		lg.Fatal("No websocket url found.")
	}
	c, _, err := websocket.DefaultDialer.Dial(ws, nil)
	if err != nil {
		lg.Fatal(err)
	}
	return c
}

// UpdateDOMEvent takes the event and, for a certain subset of events, makes sure that the current DOM object is updated.
func UpdateDOMEvent(frame *Frame, method string, event json.Unmarshaler) {
	switch method {
	case dom.EventDOMDocumentUpdated:
		frame.DOM = nil
	case dom.EventDOMSetChildNodes:
		frame.setChildNodes(&event.(*dom.SetChildNodesReply).Nodes)
	}
}

// Read reads replies from the server over the websocket.
func Read(frame *Frame) {
	for {
		_, message, err := frame.Conn.ReadMessage()
		if err != nil {
			frame.Browser.Log.Println("Read error:", err)
			return
		}
		if frame.LogLevel > LogBasic {
			frame.Browser.Log.Printf(".RAW: %s\n", message)
		}

		m := Message{}
		err = json.Unmarshal(message, &m)
		if err != nil {
			frame.Browser.Log.Fatal("Unmarshal error:", err)
		}
		//frame.Browser.Log.Printf(".DEC: %+v\n", m)

		if m.Method == "Runtime.consoleAPICalled" {
			frame.Browser.Console.Print(string(message))
		}

		hasCommand, hasEvent := false, false
		if hasCommand = frame.HasCommandID(m.ID); hasCommand {
			// All messages with an ID matching a command are set here.
			err := frame.SetResult(frame, m)
			if err != nil {
				// An unmarshal error means that the server sent an error message.  Retry.
				frame.ActionChan <- frame.ToJSON()
				continue
			}
		} else if hasEvent = frame.HasEvent(m.Method); hasEvent {
			// Check and then set Events related to the current Action.
			err = frame.SetEvent(frame, m.Method, m)
			if err != nil {
				frame.Browser.Log.Fatal(err)
			}
		}

		// If matched a command or an event, then this message is fully processed.
		if hasCommand || hasEvent {
			if frame.IsComplete() {
				frame.Browser.Log.Printf("Action Completed %s %s", frame.GetCommandMethod(), frame.GetFrameID())
				frame.Clear()
				frame.CacheCompleteChan <- struct{}{}
			} else if !frame.IsCommandComplete() {
				frame.Browser.Log.Printf("Action Next Command %s %s", frame.GetCommandMethod(), frame.GetFrameID())
				frame.ActionChan <- frame.ToJSON()
				frame.CommandChan <- frame.CommandTimeout()
			} else {
				frame.Browser.Log.Printf("Action Event Waiting %s %s", frame.GetCommandMethod(), frame.GetFrameID())
			}
			continue
		}

		// Generic unmarshaler for all other Events.
		e, ok := lib.GetEventUnmarshaler(m.Method)
		if ok {
			if len(m.Result) > 0 {
				err := e.UnmarshalJSON(m.Result)
				if err != nil {
					frame.Browser.Log.Fatal("Unmarshal error:", err, m.Result)
				}
			}
			if len(m.Params) > 0 {
				err := e.UnmarshalJSON(m.Params)
				if err != nil {
					frame.Browser.Log.Fatal("Unmarshal error:", err, m.Params)
				}
			}
			if frame.LogLevel > LogBasic {
				frame.Browser.Log.Printf(".GOT event %+v\n", e)
			}
			UpdateDOMEvent(frame, m.Method, e)
		} else {
			if frame.LogLevel > LogBasic {
				frame.Browser.Log.Printf(".SKP event %s %s %s\n", m.Method, m.Params, m.Result)
			}
		}
	}
}

// Write writes requests to the server over the websocket.
func Write(frame *Frame) {
	osInterrupt := make(chan os.Signal, 1)
	signal.Notify(osInterrupt, os.Interrupt)

	for {
		select {
		case command := <-frame.ActionChan:
			frame.Browser.Log.Printf("!REQ: %s\n", command)
			err := frame.Conn.WriteMessage(websocket.TextMessage, command)
			if err != nil {
				frame.Browser.Log.Println("write:", err)
				return
			}
		case <-frame.AllComplete:
			SendClose(frame.Browser.Log, frame.Conn)
			return
		case <-osInterrupt:
			SendClose(frame.Browser.Log, frame.Conn)
			return
		}
	}
}

// SendClose closes the websocket.
func SendClose(lg *log.Logger, c *websocket.Conn) {
	// Cleanly close the connection by sending a close message and then waiting (with timeout) for the server to close the connection.
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		lg.Println("write close err:", err)
	}
}
