package main

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/vinz1911/fastsocketserver"
)

func main() {
	var socket = fastsocketserver.FastSocketServer{}
	socket.onBinaryMessage = func(str string) {

	}
	var upgrader = websocket.Upgrader{}
	log.Println("Hello World")
}
