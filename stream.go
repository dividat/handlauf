package main

import (
	"encoding/json"
	"sync/atomic"

	"github.com/256dpi/gomqtt/broker"
	"github.com/256dpi/gomqtt/packet"
	"github.com/256dpi/gomqtt/transport"
)

type stream struct {
	clients int64
	backend broker.Backend
	engine  *broker.Engine
	server  transport.Server
}

func createStream(addr string) (*stream, error) {
	// prepare backend
	backend := broker.NewMemoryBackend()

	// prepare engine
	engine := broker.NewEngine(backend)

	// launch server
	server, err := transport.Launch(addr)
	if err != nil {
		return nil, err
	}

	// accept
	engine.Accept(server)

	// prepare stream
	stream := &stream{
		backend: backend,
		engine:  engine,
		server:  server,
	}

	// set logger
	backend.Logger = func(event broker.LogEvent, client *broker.Client, generic packet.Generic, message *packet.Message, err error) {
		switch event {
		case broker.NewConnection:
			atomic.AddInt64(&stream.clients, 1)
		case broker.LostConnection:
			atomic.AddInt64(&stream.clients, -1)
		}
	}

	return stream, nil
}

func (s *stream) emit(values sample) {
	// encode
	buf, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	// broadcast
	_ = s.backend.Publish(nil, &packet.Message{
		Topic:   "values",
		Payload: buf,
	}, nil)
}

func (s *stream) close() {
	_ = s.server.Close()
	s.engine.Close()
}
