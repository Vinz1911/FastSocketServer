package main

// Example of a `speed test server` implementation with
// the help of the FastSocket Protocol

import (
	"net"
	"strconv"

	"./fastsocket"
	"github.com/fatih/color"
)

// Mapper is a Helper to convert data types
type Mapper struct{}

func main() {
	printError := color.New(color.FgRed).PrintlnFunc()
	printInfo := color.New(color.FgYellow).PrintlnFunc()

	printInfo("[INFO]: creating server...")
	port := uint16(7878)
	mapper := Mapper{}
	server := fastsocket.Server{}
	server.OnReady = func(socket net.Conn) {
		printInfo("[INFO]: new connection id:", socket)
		printInfo("[INFO]: local address:", socket.LocalAddr())
		printInfo("[INFO]: remote address", socket.RemoteAddr())
	}
	// respond on binary message
	server.OnBinaryMessage = func(data []byte, socket net.Conn) {
		size := len(data)
		message := []byte(mapper.intToStr(size))
		server.SendMessage(fastsocket.TextMessage, message, socket)

	}
	// respond on text message
	server.OnTextMessage = func(str string, socket net.Conn) {
		response := str
		size := mapper.strToInt(response)
		var buffer []byte
		if size <= 0 {
			buffer = make([]byte, 1)
		} else {
			buffer = make([]byte, size)
		}
		server.SendMessage(fastsocket.BinaryMessage, buffer, socket)
	}
	// respond on error
	server.OnError = func(err error, socket net.Conn) {
		printError("[ERROR]: ", err)
	}
	// respond on close
	server.OnClose = func(socket net.Conn) {
		printInfo("[INFO]: connection closed and removed id:", socket)
	}
	printInfo("[INFO]: server started on port:", port)

	err := server.Start(fastsocket.TCPTransfer, port)
	if err != nil {
		printError("[ERROR]: ", err)
		return
	}
}

// convert a string to an integer
func (*Mapper) strToInt(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
	}
	return value
}

// convert an integer to a string
func (*Mapper) intToStr(integer int) string {
	value := strconv.Itoa(integer)
	return value
}
