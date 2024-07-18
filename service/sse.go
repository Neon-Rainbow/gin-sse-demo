package service

import (
	"SSE/model"
	"log"
)

type ClientChan chan model.Message

type Client struct {
	name string
	ch   ClientChan
}

type SSEvent struct {
	NewClient   chan Client
	CloseClient chan Client
	Message     chan model.Message
	Clients     map[string]ClientChan
	Messages    map[string][]model.Message
}

func NewSSEvent() *SSEvent {
	sse := &SSEvent{
		NewClient:   make(chan Client),
		CloseClient: make(chan Client),
		Message:     make(chan model.Message),
		Clients:     make(map[string]ClientChan),
		Messages:    make(map[string][]model.Message),
	}
	go sse.listen()
	return sse
}

func (sse *SSEvent) listen() {
	for {
		select {
		case client := <-sse.NewClient:
			sse.Clients[client.name] = client.ch
			for k, v := range sse.Clients {
				if k != client.name {
					v <- model.Message{
						Kind: "online",
						From: client.name,
					}
					client.ch <- model.Message{
						Kind: "online",
						From: k,
					}
				}
			}
			if message, ok := sse.Messages[client.name]; ok && len(message) > 0 {
				for _, v := range message {
					client.ch <- v
				}
				delete(sse.Messages, client.name)
				log.Printf("send offline message to %s", client.name)
			}
		case client := <-sse.CloseClient:
			if c, ok := sse.Clients[client.name]; ok {
				close(c)
				delete(sse.Clients, client.name)
				for k, v := range sse.Clients {
					if k != client.name {
						v <- model.Message{
							Kind: "offline",
							From: client.name,
						}
					}
				}
			}
		case msg := <-sse.Message:
			if c, ok := sse.Clients[msg.To]; ok {
				c <- msg
			} else {
				if messages, ok := sse.Messages[msg.To]; !ok {
					sse.Messages[msg.To] = append(messages, msg)
				} else {
					sse.Messages[msg.To] = append(sse.Messages[msg.To], msg)
				}
			}
		}
	}
}
