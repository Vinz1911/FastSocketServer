package fastsocket

import (
	"errors"
	"log"
	"net"
	"strconv"
)

// Server represents the implementation of
// the FastSocket Protocol (Server sided)
type Server struct {
	socket net.Listener
	// Closure for incoming text messages
	OnTextMessage StringClosureSocket
	// Closure for incoming data messages
	OnBinaryMessage ByteClosureSocket
}

// Frame is a struct to Create and parse the
// Protocol Messages
// Custom TCP Communication Protocol Framing
// 0                 1                              N                 N
// +-----------------+------------------------------+-----------------+
// |0 1 2 3 4 5 6 7 8|        ... Continue          |0 1 2 3 4 5 6 7 8|
// +-----------------+------------------------------+-----------------+
// |   O P C O D E   |         Payload Data...      |  F I N B Y T E  |
// +-----------------+------------------------------+-----------------+
//
type Frame struct {
	cache    []byte
	onText   ByteClosure
	onBinary ByteClosure
}

// Mapper is a Helper to convert Datatypes
type Mapper struct{}

// ByteClosureSocket Descripes a clsoure which returns a byte array
type ByteClosureSocket func([]byte, net.Conn)

// StringClosureSocket Descripes a clsoure which returns a string
type StringClosureSocket func(string, net.Conn)

// ByteClosure Descripes a clsoure which returns a byte array
type ByteClosure func([]byte)

// StringClosure Descripes a clsoure which returns a string
type StringClosure func(string)

// ControlCode are the Codes which
// describe the fragmented frame
type ControlCode uint8

// Opcode are Status Codes for
// the Message
type Opcode uint8

const (
	// continueByte is a placeholder
	continueByte ControlCode = 0x0
	// text holds the byte for a text message
	text Opcode = 0x1
	// binary holds the byte for a binary message
	binary Opcode = 0x2
	// finByte holds the ControlCode for `end of a message`
	finByte ControlCode = 0x03
	// acceptByte is for ACK if the handshake succeed
	acceptByte ControlCode = 0x06
	// connectionClose holds the byte for convenient closed connections (not used now)
	connectionClose Opcode = 0x8
	// socketID is the unique ID to identify protocol conformance connections (used for the handshake)
	socketID string = "6D8EDFD9-541C-4391-9171-AD519876B32E"
	// maximumLength is the maximum buffer read length
	maximumLength int = 8192
)

// Start starts the FastSocketServer and handles all incoming connection
func (server *Server) Start(port uint16) {
	mapper := Mapper{}
	err := errors.New("")
	server.socket, err = net.Listen("tcp", ":"+mapper.intToStr(int(port)))
	if err != nil {
		log.Println(err)
	}
	log.Println("FastSocket Server started on Port:", mapper.intToStr(int(port)))
	defer server.socket.Close()
	for {
		connection, err := server.socket.Accept()
		if err != nil {
			log.Panicln(err)
		}
		go server.loop(connection)
	}
}

// Stop closes all tcp connections
func (server *Server) Stop() {
	server.socket.Close()
}

func (server *Server) loop(socket net.Conn) {
	locked := false
	frame := Frame{}
	server.frameClosures(&frame, socket)
	for {
		buffer := make([]byte, maximumLength)
		size, err := socket.Read(buffer)
		if err != nil {
			return
		}
		data := buffer[:size]
		switch locked {
		case true:
			frame.parse(&data)
		case false:
			if string(data) == socketID {
				locked = true
				socket.Write([]byte{byte(acceptByte)})
			} else {
				socket.Close()
			}
		}
	}
}

// Handle Response for Speedtest
func (server *Server) frameClosures(frame *Frame, socket net.Conn) {
	frame.onText = func(data []byte) {
		message := string(data)
		server.OnTextMessage(message, socket)
	}
	frame.onBinary = func(data []byte) {
		server.OnBinaryMessage(data, socket)
	}
}

// SendData is used to send a binary message to the client
func (*Server) SendData(data []byte, socket net.Conn) {
	frame := Frame{}
	byted := data
	message := frame.create(&byted, binary)
	socket.Write(*message)
}

// SendString is used to send a text message to the client
func (*Server) SendString(str string, socket net.Conn) {
	frame := Frame{}
	byted := []byte(str)
	message := frame.create(&byted, text)
	socket.Write(*message)
}

// Create a FastSocket Protocol compliant frame
// this functions add the neccessary bytes to the buffer
func (*Frame) create(data *[]byte, opcode Opcode) *[]byte {
	buffer := []byte{}
	buffer = append(buffer, byte(opcode))
	buffer = append(buffer, *data...)
	buffer = append(buffer, byte(finByte))
	return &buffer
}

// Parse received Data into a FastSocket compliant
// frame/message
func (frame *Frame) parse(data *[]byte) {
	if len(*data) <= 0 {
		return
	}
	frame.cache = append(frame.cache, *data...)
	if (*data)[len(*data)-1] == byte(finByte) {
		if frame.cache[0] == byte(text) {
			var message = frame.trimmedFrame(frame.cache)
			frame.onText(message)
		}
		if frame.cache[0] == byte(binary) {
			var message = frame.trimmedFrame(frame.cache)
			frame.onBinary(message)
		}
		frame.cache = []byte{}
	}
}

// helper function to trimm a frame into a message
func (*Frame) trimmedFrame(data []byte) []byte {
	var trimmed = data
	if len(data) >= 3 {
		_, trimmed = trimmed[0], trimmed[1:]
		trimmed = trimmed[:len(trimmed)-1]
	} else {
		log.Println(errors.New("Cannot parse Frame"))
	}
	return trimmed
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
