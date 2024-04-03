package Broadcaster

type Broadcaster struct {
	subscribers []chan interface{}
	broadcast   chan interface{}
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		subscribers: make([]chan interface{}, 0),
		broadcast:   make(chan interface{}),
	}
}

func (b *Broadcaster) Subscribe(bufferSize int) chan interface{} {
	ch := make(chan interface{}, bufferSize)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *Broadcaster) Broadcast(value interface{}) {
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
