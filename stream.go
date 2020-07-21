package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/googollee/go-socket.io"
)

type stream struct {
	sock *socketio.Server
	http *http.Server
}

func createStream() (*stream, error) {
	// create server
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	// on connect
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	// on event
	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	// on error
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	// in disconnect
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	// serve
	go func() {
		panic(server.Serve().Error())
	}()

	return &stream{
		sock: server,
	}, nil
}

func (s *stream) run(addr string) error {
	// prepare errors
	errs := make(chan error, 2)

	// serve connections
	go func() {
		errs <- s.sock.Serve()
	}()

	// accept connections
	go func() {
		s.http = &http.Server{Addr: addr, Handler: s.sock}
		errs <- s.http.ListenAndServe()
	}()

	// await error
	err := <-errs

	// close
	s.close()

	return err
}

func (s *stream) emit(values sample) {
	// encode
	buf, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	// broadcast
	s.sock.BroadcastToRoom("main", "values", string(buf))
}

func (s *stream) close() {
	_ = s.sock.Close()
	_ = s.http.Close()
}
