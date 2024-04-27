package server

import (
	"log/slog"
	"net"
	"time"
)

type Server struct {
	ln         net.Listener
	listenAddr string
	quitch     chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
	}
}
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptConn()
	<-s.quitch
	return nil
}

func (s *Server) acceptConn() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("Error accepting connection: ", err)
			continue
		}
		slog.Info("Connected to: ", conn.RemoteAddr())
		go s.writeStream(conn)
	}
}

func (s *Server) writeStream(conn net.Conn) {
	for {
		conn.Write([]byte(("Pinging Connection every 3 seconds...\n")))
		time.Sleep(3 * time.Second)
	}
}
