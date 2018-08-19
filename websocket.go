package cdp

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// GetWebsocket returns a websocket connection to the running browser.
func GetWebsocket() *websocket.Conn {
	r, err := http.Get("http://localhost:9222/json")
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

		if ac.HasStepID(m.ID) {
			err := ac.SetResult(m)
			if err != nil {
				log.Fatal(err)
			}
			stepComplete <- struct{}{}
		} else {
			// Check for events related to the current Action
			if ac.HasEvent(m.Method) {
				pi, err := UnmarshalIds(m)
				if err != nil {
					log.Fatal("Unmarshal error:", err)
				}
				log.Printf(".IDS: %+v\n", pi)

				err = ac.SetEvent(m.Method, m, pi)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Printf("SKIP event %s %s %s\n", m.Method, m.Params, m.Result)
			}
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
