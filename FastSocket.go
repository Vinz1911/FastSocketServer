package fastsocket

import (
	"errors"
	"log"
	"net"
	"strconv"
)

// byte closure describe a closure which returns a byte array
type byteClosure func([]byte, net.Conn)

// string closure describe a closure which returns a string
type stringClosure func(string, net.Conn)

// operational codes are used to
// control the framing, handles handshake and more
type operationalCode uint8

// the message type for framing
type messageType uint8

const (
	// continueByte is a placeholder `UNUSED`
	// continueByte operationalCode = 0x0
	// text holds the byte for a text message
	TextMessage messageType = 0x1

	// binary holds the byte for a binary message
	BinaryMessage messageType = 0x2

	// finByte holds the ControlCode for `end of a message`
	finByte operationalCode = 0x03

	// acceptByte is for ACK if the handshake succeed
	acceptByte operationalCode = 0x06

	// connectionClose holds the byte for convenient closed connections `UNUSED`
	// connectionClose operationalCode = 0x8

	// socketID is the unique ID to identify protocol conformance connections (used for the handshake)
	socketID string = "6D8EDFD9-541C-4391-9171-AD519876B32E"

	// maximumLength is the maximum buffer read length
	maximumLength int = 8192

	// maximum frame size
	maximumContentLength int = 16777216
)

// +---------------------------+
// |        S E R V E R        |
// +---------------------------+
// Server represents the implementation of
// the FastSocket Protocol (Server sided)
type Server struct {
	// transfer
	transfer transfer
	// Closure for incoming text messages
	OnTextMessage stringClosure
	// Closure for incoming data messages
	OnBinaryMessage byteClosure
}
// Start starts the FastSocketServer and handles all incoming connection
func (server *Server) Start(port uint16) {
	server.transfer = transfer{}
	server.transferClosure()
	server.transfer.start(port)
}
// Stop closes all tcp connections
func (server *Server) Stop() {
	server.transfer.stop()
}
// SendMessage is used to send data or string based messages to the client
func (server *Server) SendMessage(messageType messageType, data *[]byte, conn net.Conn) {
	frame := frame{}
	if messageType == TextMessage {
		message, err := frame.create(data, TextMessage)
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		server.write(message, conn)
	}
	if messageType == BinaryMessage {
		message, err := frame.create(data, BinaryMessage)
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		server.write(message, conn)
	}
}
// internal used transfer read loop
func (server *Server) transferClosure() {
	server.transfer.onConnection = func(conn net.Conn) {
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
}
// internal sending function
func (server *Server) write(data *[]byte, conn net.Conn) {
	_, err := conn.Write(*data)
	if err != nil {
		log.Println(err)
		return
	}
}
// Handle Response for speed test
func (server *Server) frameClosures(frame *frame, conn net.Conn) {
	frame.onText = func(data []byte) {
		message := string(data)
		server.OnTextMessage(message, conn)
	}
	frame.onBinary = func(data []byte) {
		server.OnBinaryMessage(data, conn)
	}
}
// +-----------------------------+
// |        F R A M I N G        |
// +-----------------------------+
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
	onText   func([]byte)
	onBinary func([]byte)
}
// Create a FastSocket Protocol compliant frame
// this functions add the necessary bytes to the buffer
func (*frame) create(data *[]byte, messageType messageType) (*[]byte, error) {
	var buffer []byte
	buffer = append(buffer, byte(messageType))
	buffer = append(buffer, *data...)
	buffer = append(buffer, byte(finByte))
	if len(buffer) > maximumContentLength {
		return nil, errors.New("write buffer overflow")
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
		return errors.New("read buffer overflow")
	}
	if (*data)[len(*data)-1] == byte(finByte) {
		if frame.cache[0] == byte(TextMessage) {
			message, err := frame.trimmedFrame(frame.cache)
			if err != nil {
				return err
			}
			frame.onText(message)
		}
		if frame.cache[0] == byte(BinaryMessage) {
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
// helper function to trim a frame into a message
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
// +-------------------------------+
// |        T R A N S F E R        |
// +-------------------------------+
// transfer is a raw tcp socket
type transfer struct {
	// the tcp listener
	listener net.Listener
	// go routine closure with the connection
	onConnection func(net.Conn)
}
// start the network listener
func (transfer *transfer) start(port uint16) {
	var err error
	transfer.listener, err = net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		log.Println(err)
	}
	log.Println("FastSocket Server started on Port:", strconv.Itoa(int(port)))
	defer transfer.listener.Close()
	for {
		connection, err := transfer.listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go transfer.onConnection(connection)
	}
}
// invalidate all current running tcp connections
func (transfer *transfer) stop() {
	transfer.listener.Close()
}