package pubsub

type Topic struct {
    name       string
    subscribers map[*Client]struct{}
    register    chan *Client
    unregister  chan *Client
    broadcast   chan []byte
}

func (t *Topic) Register(cl *Client){
	t.register <- cl
}

func (t *Topic) Unregister(cl *Client){
	t.unregister <- cl
}

func (t *Topic) NewMessage(b []byte){
	t.broadcast <- b
}

func newTopic(name string)*Topic{
	return &Topic{
		name: name,
		subscribers: map[*Client]struct{}{},
		register: make(chan *Client),
		unregister: make(chan *Client),
		broadcast: make(chan []byte),
	}
}

func (t *Topic) serve() {
    for {
        select {
        case cl := <-t.register:
            t.subscribers[cl] = struct{}{}

        case cl := <-t.unregister:
            delete(t.subscribers, cl)

        case msg := <-t.broadcast:
            for cl := range t.subscribers {
                select {
                case cl.input <- msg:
                default:
					panic("slow client")
                    // delete(t.subscribers, cl)
                }
            }
        }
    }
}

