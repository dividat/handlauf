package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

var clients sync.Map
var numClients int64
var server *http.Server

func stream(addr string) {
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
			fmt.Printf("stream: %s\n", err.Error())
			return
		}

		// add client
		fmt.Printf("stream: added client %s\n", conn.RemoteAddr().String())
		clients.Store(conn, nil)
		atomic.AddInt64(&numClients, 1)
	})

	// run server
	server = &http.Server{Addr: addr, Handler: handler}
	panic(server.ListenAndServe())
}

func emit(values sample) {
	// build payload
	payload := make([]byte, 0, 256)
	for i, v := range values {
		if i > 0 {
			payload = append(payload, ',')
		}
		payload = strconv.AppendFloat(payload, v, 'f', 2, 64)
	}

	// broadcast
	clients.Range(func(key, _ interface{}) bool {
		// get conn
		conn := key.(*websocket.Conn)

		// write message
		err := conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			fmt.Printf("emit: %s\n", err.Error())
			_ = conn.Close()
			clients.Delete(key)
			atomic.AddInt64(&numClients, -1)
			fmt.Printf("emit: removed client %s", conn.RemoteAddr().String())
		}

		return true
	})
}
