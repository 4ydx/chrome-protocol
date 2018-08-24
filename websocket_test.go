package cdp

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestShutdown(t *testing.T) {
	srv := Serve()
	defer ServerClose(srv)

	c := GetWebsocket(8080)

	// Direct use of the connection to see that data is sent/received.
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
