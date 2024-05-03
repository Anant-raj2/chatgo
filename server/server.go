package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
)

type Message struct {
	addr string
	msg  []byte
}

type Server struct {
	ln         net.Listener
	listenAddr string
	quitch     chan struct{}

	peers      map[net.Conn]bool
	peersMutex sync.RWMutex

	msgch chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		peers:      make(map[net.Conn]bool),
		peersMutex: sync.RWMutex{},
		msgch:      make(chan Message),
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
		fmt.Println("Connected to: ", conn.RemoteAddr().String())
		s.peersMutex.Lock()
		s.peers[conn] = true
		s.peersMutex.Unlock()
		go s.readStream(conn)
		go func() {
			for {
				msg := <-s.msgch
				s.broadcastStream(msg.addr, msg.msg)
			}
		}()
	}
}

func (s *Server) readStream(conn net.Conn) {
	defer func() {
		s.peersMutex.Lock()
		s.peers[conn] = false
		s.peersMutex.Unlock()
		conn.Close()
	}()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println(conn.RemoteAddr().String(), " was disconnected.")
				s.broadcastStream(conn.RemoteAddr().String(), []byte("Disconnected."))
				return
			}
			slog.Error("Error reading stream: ", err)
		}
		s.msgch <- Message{
			addr: conn.RemoteAddr().String(),
			msg:  buf[:n],
		}
	}
}

func (s *Server) broadcastStream(addr string, msg []byte) {
	for index, val := range s.peers {
		if val {
			_, writeErr := index.Write([]byte(fmt.Sprintf("%s: %s", addr, msg)))
			if writeErr != nil {
				slog.Error("There is an error writing: ", writeErr)
				return
			}
		}
	}
	fmt.Println("Message Broadcasted")
}

func main() {
	var server *Server = NewServer(":3000")
	panic(server.Start())
}
