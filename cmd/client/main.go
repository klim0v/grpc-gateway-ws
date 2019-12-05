package main

import (
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/url"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8000", Path: "/ws/echo"}
	log.Printf("connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		if err == websocket.ErrBadHandshake {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
				return
			}
			log.Fatalf("handshake failed with status %d, err: %s", resp.StatusCode, body)
		} else {
			log.Fatal("dial:", err)
		}
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s	", message)
		}
	}()

	<-done
}
