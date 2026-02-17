package server

import (
	"andrew_chat/intenal/config"
	"andrew_chat/intenal/domain"
	"net"
	"time"

	"github.com/google/uuid"
)


const (
	StatusConnected = iota
	StatusDisconnected
	StatusConnecting
)

type ServerService struct {
	server domain.Server
	conn   net.Conn
	done   chan struct{}
}

func NewServerService() *ServerService {
	return &ServerService{}
}

func (ss *ServerService) Connect(srv domain.Server) error {
	// ss.server = srv
	// conn, err := net.Dial("udp", srv.Addr)
	// if err != nil {
	// 	return err
	// }
	// ss.conn = conn
	return nil
}

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

func (ss *ServerService) GetServers() []domain.Server{
	return config.GetServers()
}