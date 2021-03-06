package main

// Example of a `speed test server` implementation with
// the help of the FastSocket Protocol
import (
	"./fastsocket"
	"errors"
	"flag"
	"github.com/fatih/color"
	"math"
	"net"
	"strconv"
)

// Mapper is a Helper to convert data types
type Mapper struct{}

func main() {
	printError := color.New(color.FgRed).PrintlnFunc()
	printInfo := color.New(color.FgCyan).PrintlnFunc()

	var port uint
	flag.UintVar(&port, "port", 7878, "server port")
	flag.Parse()
	if port > math.MaxUint16 {
		printError(errors.New("value overflows uint16, please check given port"))
		return
	}
	mapper := Mapper{}
	server := fastsocket.Server{}
	server.OnReady = func(socket net.Conn) {
		printInfo("[INFO]: new connection id:", socket)
		printInfo("[INFO]: local address:", socket.LocalAddr())
		printInfo("[INFO]: remote address", socket.RemoteAddr())
	}
	// respond on binary message
	server.OnDataMessage = func(socket net.Conn, data []byte) {
		size := len(data)
		message := mapper.intToStr(size)
		server.SendStringMessage(socket, message)
	}
	// respond on text message
	server.OnStringMessage = func(socket net.Conn, str string) {
		response := str
		size := mapper.strToInt(response)
		var message []byte
		if size <= 0 {
			message = make([]byte, 1)
		} else {
			message = make([]byte, size)
		}
		server.SendDataMessage(socket, message)
	}
	// respond on error
	server.OnError = func(socket net.Conn, err error) {
		printError("[ERROR]: ", err)
	}
	// respond on close
	server.OnClose = func(socket net.Conn) {
		printInfo("[INFO]: connection closed and removed id:", socket)
	}
	printInfo("[INFO]: creating server...")
	printInfo("[INFO]: server started on port:", port)
	err := server.Start(fastsocket.TCPTransfer, uint16(port))
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
