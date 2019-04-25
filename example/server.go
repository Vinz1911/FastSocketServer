package main

import (
	"log"
	"net"

	"github.com/Vinz1911/fastsocket"
	//"github.com/vinz1911/fastsocket"
)

func main() {
	var socket = fastsocket.Server{}
	socket.OnBinaryMessage = func(str []byte, socket net.Conn) {
	}
	socket.Start(3333)
	log.Println("Hello World")
}
