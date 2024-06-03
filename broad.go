package main

import "github.com/vtb-link/bianka/proto"

type Broadcaster struct {
	subscribers []chan *proto.CmdDanmuData
	broadcast   chan *proto.CmdDanmuData
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		subscribers: make([]chan *proto.CmdDanmuData, 0),
		broadcast:   make(chan *proto.CmdDanmuData),
	}
}

func (b *Broadcaster) Subscribe(bufferSize int) chan *proto.CmdDanmuData {
	ch := make(chan *proto.CmdDanmuData, bufferSize)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *Broadcaster) Broadcast(value *proto.CmdDanmuData) {
	b.broadcast <- value
}

func (b *Broadcaster) Start() {
	go func() {
		for {
			value := <-b.broadcast
			for _, subscriber := range b.subscribers {
				subscriber <- value
			}
		}
	}()
}
