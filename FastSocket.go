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
	// socket
	socket net.Listener
	// Closure for incoming text messages
	OnTextMessage stringClosureSocket
	// Closure for incoming data messages
	OnBinaryMessage byteClosureSocket
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
type frame struct {
	cache    []byte
	onText   byteClosure
	onBinary byteClosure
}

// ByteClosureSocket Descripes a closure which returns a byte array
type byteClosureSocket func([]byte, net.Conn)

// StringClosureSocket Descripes a closure which returns a string
type stringClosureSocket func(string, net.Conn)

// ByteClosure Descripes a closure which returns a byte array
type byteClosure func([]byte)

// StringClosure Descripes a closure which returns a string
type stringClosure func(string)

// ControlCode are the Codes which
// describe the fragmented frame
type controlCode uint8

// Opcode are Status Codes for
// the Message
type opcode uint8

const (
	// continueByte is a placeholder
	continueByte controlCode = 0x0
	// text holds the byte for a text message
	text opcode = 0x1
	// binary holds the byte for a binary message
	binary opcode = 0x2
	// finByte holds the ControlCode for `end of a message`
	finByte controlCode = 0x03
	// acceptByte is for ACK if the handshake succeed
	acceptByte controlCode = 0x06
	// connectionClose holds the byte for convenient closed connections (not used now)
	connectionClose opcode = 0x8
	// socketID is the unique ID to identify protocol conformance connections (used for the handshake)
	socketID string = "6D8EDFD9-541C-4391-9171-AD519876B32E"
	// maximumLength is the maximum buffer read length
	maximumLength int = 8192
	// maximum frame size
	maximumContentLength int = 16777216
)

// Start starts the FastSocketServer and handles all incoming connection
func (server *Server) Start(port uint16) {
	err := errors.New("")
	server.socket, err = net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		log.Println(err)
	}
	log.Println("FastSocket Server started on Port:", strconv.Itoa(int(port)))
	defer server.socket.Close()
	for {
		connection, err := server.socket.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go server.loop(connection)
	}
}

// Stop closes all tcp connections
func (server *Server) Stop() {
	server.socket.Close()
}

// SendData is used to send a binary message to the client
func (server *Server) SendData(data []byte, conn net.Conn) {
	frame := frame{}
	message, err := frame.create(&data, binary)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	server.write(message, conn)
}

// SendString is used to send a text message to the client
func (server *Server) SendString(str string, conn net.Conn) {
	frame := frame{}
	data := []byte(str)
	message, err := frame.create(&data, text)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	server.write(message, conn)
}

func (server *Server) loop(conn net.Conn) {
	mutexLock := false
	frame := frame{}
	server.frameClosures(&frame, conn)
	for {
		buffer := make([]byte, maximumLength)
		size, err := conn.Read(buffer)
		if err != nil {
			return
		}
		data := buffer[:size]
		switch mutexLock {
		case true:
			err := frame.parse(&data)
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}
		case false:
			if string(data) == socketID {
				mutexLock = true
				conn.Write([]byte{byte(acceptByte)})
			} else {
				conn.Close()
			}
		}
	}
}

// internal sending function
func (server *Server) write(data *[]byte, conn net.Conn) {
	_, err := conn.Write(*data)
	if err != nil {
		log.Println(err)
		return
	}
}

// Handle Response for Speedtest
func (server *Server) frameClosures(frame *frame, conn net.Conn) {
	frame.onText = func(data []byte) {
		message := string(data)
		server.OnTextMessage(message, conn)
	}
	frame.onBinary = func(data []byte) {
		server.OnBinaryMessage(data, conn)
	}
}

// Create a FastSocket Protocol compliant frame
// this functions add the neccessary bytes to the buffer
func (*frame) create(data *[]byte, op opcode) (*[]byte, error) {
	buffer := []byte{}
	buffer = append(buffer, byte(op))
	buffer = append(buffer, *data...)
	buffer = append(buffer, byte(finByte))
	if len(buffer) > maximumContentLength {
		return nil, errors.New("writebuffer overflow")
	}
	return &buffer, nil
}

// Parse received Data into a FastSocket compliant
// frame/message
func (frame *frame) parse(data *[]byte) error {
	if len(*data) <= 0 {
		return errors.New("zero data")
	}
	frame.cache = append(frame.cache, *data...)
	if len(frame.cache) > maximumContentLength {
		return errors.New("readbuffer overflow")
	}
	if (*data)[len(*data)-1] == byte(finByte) {
		if frame.cache[0] == byte(text) {
			message, err := frame.trimmedFrame(frame.cache)
			if err != nil {
				return err
			}
			frame.onText(message)
		}
		if frame.cache[0] == byte(binary) {
			message, err := frame.trimmedFrame(frame.cache)
			if err != nil {
				return err
			}
			frame.onBinary(message)
		}
		frame.cache = []byte{}
	}
	return nil
}

// helper function to trimm a frame into a message
func (*frame) trimmedFrame(data []byte) ([]byte, error) {
	var trimmed = data
	if len(data) >= 3 {
		_, trimmed = trimmed[0], trimmed[1:]
		trimmed = trimmed[:len(trimmed)-1]
	} else {
		return nil, errors.New("cannot parse Frame")
	}
	return trimmed, nil
}
