package main

// Example of a `Speedtest Server` implementation with
// the help of the FastSocket Protocol

import (
	"net"
	"strconv"

	"github.com/vinz1911/fastsocketserver"
)

// Mapper is a Helper to convert datatypes
type Mapper struct{}

func main() {
	port := uint16(8080)
	mapper := Mapper{}
	server := fastsocket.Server{}
	server.OnBinaryMessage = func(data []byte, socket net.Conn) {
		size := len(data)
		message := mapper.intToStr(size)
		server.SendString(message, socket)
	}
	server.OnTextMessage = func(str string, socket net.Conn) {
		size := str
		message := make([]byte, mapper.strToInt(size))
		server.SendData(message, socket)
	}
	server.Start(port)
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
