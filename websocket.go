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
	"time"
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
		fmt.Println("-------------")
		entry := arr[i].(map[string]interface{})
		for k, v := range entry {
			fmt.Printf("%s %+v\n", k, v)
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
func Read(c *websocket.Conn, stepComplete chan<- struct{}, ac *ActionCache, shutdown chan<- struct{}) {
	defer func() {
		log.Println("Shutdown due to socket connection going away.")
		close(shutdown)
	}()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		log.Printf(".RAW: %s\n", message)

		m := Message{}
		err = json.Unmarshal(message, &m)
		if err != nil {
			log.Fatal("Unmarshal error:", err)
		}
		// log.Printf(".DEC: %+v\n", m)

		// All messages with an ID matching a step are set here.
		if ac.HasStepID(m.ID) {
			err := ac.SetResult(m)
			if err != nil {
				log.Fatal(err)
			}
			stepComplete <- struct{}{}
			continue
		}

		// Check and then set Events related to the current Action.
		if ac.HasEvent(m.Method) {
			err = ac.SetEvent(m.Method, m)
			if err != nil {
				log.Fatal(err)
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
			log.Printf(".GOT event %+v\n", e)
		} else {
			log.Printf(".SKP event %s %s %s\n", m.Method, m.Params, m.Result)
		}
	}
}

// Write writes requests to the server over the websocket.
func Write(c *websocket.Conn, actionChan <-chan *Action, ac *ActionCache, shutdown, allComplete <-chan struct{}) {
	osInterrupt := make(chan os.Signal, 1)
	signal.Notify(osInterrupt, os.Interrupt)

	for {
		select {
		case <-shutdown:
			log.Println("Shutdown.")
			return
		case action := <-actionChan:
			log.Printf("!REQ: %s\n", action.ToJSON())
			ac.Set(action)

			err := c.WriteMessage(websocket.TextMessage, action.ToJSON())
			if err != nil {
				fmt.Println("write:", err)
				return
			}
		case <-allComplete:
			SendClose(c, shutdown)
			return
		case <-osInterrupt:
			SendClose(c, shutdown)
			return
		}
	}
}

// SendClose closes the websocket.
func SendClose(c *websocket.Conn, shutdown <-chan struct{}) {
	// Cleanly close the connection by sending a close message and then waiting (with timeout) for the server to close the connection.
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close err:", err)
		return
	}
	select {
	case <-shutdown:
		log.Println("SendClose done.")
	case <-time.After(time.Second * 5):
		log.Println("SendClose timeout.")
	}
}
