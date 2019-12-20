package main

import (
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/url"
	"time"
)

func main() {
	c, bad := conn(url.URL{Scheme: "ws", Host: "localhost:8841", Path: "/subscribe", RawQuery: "query=tags.tx.coin%3D%27MNT%27"})
	if bad {
		return
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

	go func() {
		time.Sleep(time.Second * 10)
		_, bad := conn(url.URL{Scheme: "ws", Host: "localhost:8841", Path: "/unsubscribe", RawQuery: "query=tags.tx.coin%3D%27MNT%27"})
		if bad {
			return
		}
	}()

	defer c.Close()

	<-done
}

func conn(u url.URL) (*websocket.Conn, bool) {
	log.Printf("connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		if err == websocket.ErrBadHandshake {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
				return nil, true
			}
			log.Fatalf("handshake failed with status %d, err: %s", resp.StatusCode, body)
		} else {
			log.Fatal("dial:", err)
		}
	}
	return c, false
}
