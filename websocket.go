package cdp

import (
	"encoding/json"
	"fmt"
	"github.com/4ydx/cdp/protocol/lib"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// GetWebsocket returns a websocket connection to the running browser.
func GetWebsocket(port int) *websocket.Conn {
	r, err := http.Get(fmt.Sprintf("http://localhost:%d/json", port))
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var inter interface{}
	err = json.Unmarshal(b, &inter)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal("No websocket url found.")
	}
	c, _, err := websocket.DefaultDialer.Dial(ws, nil)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// Read reads replies from the server over the websocket.
func Read(frame *Frame) {
	for {
		_, message, err := frame.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		if frame.LogLevel > LogBasic {
			log.Printf(".RAW: %s\n", message)
		}

		m := Message{}
		err = json.Unmarshal(message, &m)
		if err != nil {
			log.Fatal("Unmarshal error:", err)
		}
		// log.Printf(".DEC: %+v\n", m)

		hasCommand, hasEvent := false, false
		if hasCommand = frame.Cache.HasCommandID(m.ID); hasCommand {
			// All messages with an ID matching a command are set here.
			err := frame.Cache.SetResult(m)
			if err != nil {
				// An unmarshal error means that the server sent an error message.  Retry.
				frame.ActionChan <- frame.Cache.ToJSON()
				continue
			}
		} else if hasEvent = frame.Cache.HasEvent(m.Method); hasEvent {
			// Check and then set Events related to the current Action.
			err = frame.Cache.SetEvent(m.Method, m)
			if err != nil {
				log.Fatal(err)
			}
		}

		// If matched a command or an event, then this message is fully processed.
		if hasCommand || hasEvent {
			if frame.Cache.IsComplete() {
				log.Printf("Action Completed %s %s", frame.Cache.GetCommandMethod(), frame.Cache.GetFrameID())
				frame.Cache.Clear()
				frame.CacheCompleteChan <- struct{}{}
			} else if !frame.Cache.IsCommandComplete() {
				log.Printf("Action Next Command %s %s", frame.Cache.GetCommandMethod(), frame.Cache.GetFrameID())
				frame.ActionChan <- frame.Cache.ToJSON()
				frame.CommandChan <- frame.Cache.CommandTimeout()
			} else {
				log.Printf("Action Event Waiting %s %s", frame.Cache.GetCommandMethod(), frame.Cache.GetFrameID())
			}
			continue
		}

		// Generic unmarshaler for all other Events.
		e, ok := lib.GetEventUnmarshaler(m.Method)
		if ok {
			if len(m.Result) > 0 {
				err := e.UnmarshalJSON(m.Result)
				if err != nil {
					log.Fatal("Unmarshal error:", err, m.Result)
				}
			}
			if len(m.Params) > 0 {
				err := e.UnmarshalJSON(m.Params)
				if err != nil {
					log.Fatal("Unmarshal error:", err, m.Params)
				}
			}
			if frame.LogLevel > LogBasic {
				log.Printf(".GOT event %+v\n", e)
			}
		} else {
			if frame.LogLevel > LogBasic {
				log.Printf(".SKP event %s %s %s\n", m.Method, m.Params, m.Result)
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
			log.Printf("!REQ: %s\n", command)
			err := frame.Conn.WriteMessage(websocket.TextMessage, command)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-frame.AllComplete:
			SendClose(frame.Conn)
			return
		case <-osInterrupt:
			SendClose(frame.Conn)
			return
		}
	}
}

// SendClose closes the websocket.
func SendClose(c *websocket.Conn) {
	// Cleanly close the connection by sending a close message and then waiting (with timeout) for the server to close the connection.
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close err:", err)
	}
}
