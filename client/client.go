package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
)

type Client struct {
	conn   net.Conn
	quitch chan struct{}
  recievech chan []byte
}

func NewClient() *Client {
	return &Client{
		quitch: make(chan struct{}),
    recievech: make(chan []byte),
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		slog.Error("Error pinging: ", err)
		return err
	}
	defer conn.Close()
	c.conn = conn
	go c.writeConn()
	go c.readConn()
	<-c.quitch
	return nil
}

func (c *Client) writeConn() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_, err := c.conn.Write([]byte(scanner.Text()))
		if err != nil {
			slog.Error("Error writing to connection: ", err)
			continue
		}
	}
}
func (c *Client) readConn() {
	buf := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			slog.Error("Error pinging: ", err)
			continue
		}
    c.recievech <- buf[:n]
	}
}
func main() {
	client := NewClient()
  go func(){
    for{
      msg:=<-client.recievech
      fmt.Println(string(msg))
    }
  }()
	panic(client.Connect())
}
