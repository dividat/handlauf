package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type stream struct {
	size    int64
	clients sync.Map
	server  *http.Server
}

func createStream(addr string) (*stream, error) {
	// prepare stream
	stream := &stream{}

	// prepare upgrader
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// prepare handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// perform upgrade
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		// add client
		stream.clients.Store(conn, nil)
		atomic.AddInt64(&stream.size, 1)
	})

	// run server
	go func() {
		stream.server = &http.Server{Addr: addr, Handler: handler}
		panic(stream.server.ListenAndServe())
	}()

	return stream, nil
}

func (s *stream) emit(values sample) {
	// encode
	payload, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	// broadcast
	s.clients.Range(func(key, _ interface{}) bool {
		// get conn
		conn := key.(*websocket.Conn)

		// write message
		err := conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			_ = conn.Close()
			s.clients.Delete(key)
			atomic.AddInt64(&s.size, -1)
			println(err.Error())
		}

		return true
	})
}

func (s *stream) close() {
	_ = s.server.Close()
}
