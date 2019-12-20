package service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{}
)

func (s *Service) Subscribe(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

	start := time.Now()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {

			time.Sleep(time.Second)
			if err := conn.WriteMessage(websocket.TextMessage, []byte("hello there!"+fmt.Sprint(time.Now().Sub(start)))); err != nil {
				_ = conn.Close()
				break
			}
		}
	}()

	<-done
}
