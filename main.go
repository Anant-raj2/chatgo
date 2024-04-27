package main

import (
  source "github.com/Anant-raj2/chatgo/server"
)

func main(){
  var server *source.Server = source.NewServer(":3000")
  panic(server.Start())
}
