package server

import (
	"andrew_chat/client/internal/config"
	"andrew_chat/client/internal/domain"

	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	StatusConnected = iota
	StatusDisconnected
	StatusConnecting
)

type ServerService struct {
	server     domain.Server
	client     http.Client
	conn       net.Conn
	done       chan struct{}
	lastUpdate time.Time
}

func NewServerService() *ServerService {
	return &ServerService{}
}

func (ss *ServerService) Connect(srv domain.Server) error {
	// if ss.conn != nil {
	// 	ss.conn.Close()
	// }

	// const attempts = 3

	// var err error
	// for i := 1; i < attempts+1; i++ {
	// 	ss.conn, err = net.Dial("udp", srv.Address)
	// 	if err != nil {
	// 		slog.Error("connection failed", slog.Any("err", err), slog.Int("attempt", i))
	// 		continue
	// 	}
	// }

	// if ss.conn == nil {
	// 	slog.Error("connection failed")
	// 	return err
	// }

	// slog.Info("connection established")
	time.Sleep(1 * time.Second)
	return nil
}

// func (ss *ServerService) CreateRoom() error {
// 	ss.client.Post()
// }

func (ss *ServerService) Terminate() {
	ss.conn.Close()
}

func (ss *ServerService) Add(server domain.Server) error {
	server.ID = uuid.NewString()

	server.CreatedAt = time.Now()
	server.UpdatedAt = time.Now()

	return config.AddServer(server)
}

func (ss *ServerService) Remove(serverID string) error {
	return config.DeleteServer(serverID)
}

func (ss *ServerService) Update(server domain.Server) error {
	server.UpdatedAt = time.Now()
	return config.UpdateServer(server)
}

func (ss *ServerService) GetServers() []domain.Server {
	return config.GetServers()
}

func (ss *ServerService) SendMessage(){}
