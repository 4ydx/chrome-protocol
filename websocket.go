package main

import (
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

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

func Read(c *websocket.Conn, stepComplete chan<- int64, as *StepCache, shutdown chan<- struct{}) {
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

		m := cdproto.Message{}
		err = m.UnmarshalJSON(message)
		if err != nil {
			log.Fatal("Unmarshal error:", err)
		}
		if step, ok := as.Get(m.ID); ok {
			err := step.Returns.UnmarshalJSON(m.Result)
			if err != nil {
				log.Fatal("Unmarshal error:", err)
			}
			stepComplete <- step.ActionId
			log.Printf(".RES: %+v\n", step)
			log.Printf("    : %+v\n", step.Params)
			log.Printf("    : %+v\n", step.Returns)
		}
	}
}

func Write(c *websocket.Conn, actions <-chan *Action, as *StepCache, shutdown, allComplete <-chan struct{}) {
	osInterrupt := make(chan os.Signal, 1)
	signal.Notify(osInterrupt, os.Interrupt)

	for {
		select {
		case <-shutdown:
			fmt.Println("shutdown")
			return
		case action := <-actions:
			as.Add(action.Step())

			log.Printf("!REQ: %s\n", action.ToJSON())
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

func SendClose(c *websocket.Conn, shutdown <-chan struct{}) {
	// Cleanly close the connection by sending a close message and then waiting (with timeout) for the server to close the connection.
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close err:", err)
		return
	}
	select {
	case <-shutdown:
		log.Println("SendClose done")
	case <-time.After(time.Second * 5):
		log.Println("SendClose timeout")
	}
}
