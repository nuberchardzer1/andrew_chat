package pubsub

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)



type PubsubServer struct {
	addr      string
	mux       *http.ServeMux
	upgrader  websocket.Upgrader
	topics   map[string]*Topic
}


func (ps *PubsubServer) servePublisher(cl *Client) {
	output := cl.Readloop()

	for b := range output {
		var msg Message
		if err := msg.Unmarshal(b); err != nil {
			slog.Error("unmarshal message failed", "err", err)
			continue
		}

		if msg.Action != MessageAction {
			slog.Warn("wrong action", "expected", MessageAction, "got", msg.Action)
			continue
		}

		topic, ok := ps.topics[msg.Topic]
		if !ok {
			slog.Warn("unknown topic", "msg", msg)
			continue
		}

		b, err := msg.Marshal()
		if err != nil {
			slog.Error("marshal message failed", "err", err)
			continue
		}
		topic.broadcast <- b
	}

	select {
	case err := <-cl.errChan:
		if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			slog.Error("publisher connection error", "err", err)
		}
		cl.Close()
	default:
		slog.Debug("publisher connection closed without error")
		cl.Close()
	}
}

func (ps *PubsubServer) serveSubscriber(cl *Client){
	cl.Writeloop()
}

func (s *PubsubServer) PublishHandler(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if _, ok := s.topics[topic]; !ok {
		http.Error(w, "unknown topic", http.StatusBadRequest)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", "err", err)
		return
	}

	client := NewClient(conn)
	slog.Info("new publisher connected", "addr", client.conn.LocalAddr())
	go s.servePublisher(client)
}

func (s *PubsubServer) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	topicNames, ok := r.URL.Query()["topic"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", "err", err)
		return
	}

	client := NewClient(conn)

	for _, topicName := range topicNames {
		if topicName == "" {
			continue
		}

		topic, ok := s.topics[topicName]
		if !ok {
			topic = newTopic(topicName)
			s.topics[topicName] = topic
			go topic.serve()
		}

		topic.Register(client)
	}

	slog.Info("new subscriber connected", "addr", client.conn.LocalAddr(), "topics", topicNames)
	go s.serveSubscriber(client)
}

func (s *PubsubServer) Run() error {
	slog.Info("PubSub server running", "addr", s.addr)
	return http.ListenAndServe(s.addr, s.mux)
}

func NewDefaultServer(addr string) *PubsubServer {
	srv := &PubsubServer{
		addr:    addr,
		mux:     http.NewServeMux(),
		topics:  make(map[string]*Topic),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}

	srv.mux.HandleFunc("/publish", srv.PublishHandler)
	srv.mux.HandleFunc("/subscribe", srv.SubscribeHandler)

	return srv
}
