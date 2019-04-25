package main

import (
	"log"
	"net"

	fastsocket "github.com/Vinz1911/fastsocketserver"
	//"github.com/vinz1911/fastsocket"
)

func main() {
	var socket = fastsocket.Server{}
	socket.OnBinaryMessage = func(str []byte, socket net.Conn) {
	}
	socket.Start(3333)
	log.Println("Hello World")
}
