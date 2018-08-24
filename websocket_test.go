package cdp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
)

var (
	addr     = "localhost:8080"
	upgrader = websocket.Upgrader{}
)

func JSON(w http.ResponseWriter, r *http.Request) {
	reply := map[string]string{
		"webSocketDebuggerUrl": fmt.Sprintf("ws://%s/ws", addr),
	}
	value := []interface{}{
		reply,
	}
	b, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-type", "application/json")
	w.Write(b)
}

func Echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	log.Println("closing websocket")
}

func Serve() *http.Server {
	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/json", JSON)
	http.HandleFunc("/ws", Echo)
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	return srv
}

func ServerClose(srv *http.Server) {
	if err := srv.Shutdown(nil); err != nil {
		panic(err)
	}
}

func TestShutdown(t *testing.T) {
	srv := Serve()
	defer ServerClose(srv)

	c := GetWebsocket(8080)

	err := c.WriteMessage(websocket.TextMessage, []byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	_, message, err := c.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if string(message) != "hello" {
		t.Fatal("expecting hello message")
	}
	SendClose(c)

	_, _, err = c.ReadMessage()
	if err == nil {
		t.Fatal("expecting the connection to be closed")
	}
	t.Log("shutdown err", err)
}
