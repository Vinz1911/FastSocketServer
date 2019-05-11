package fastsocket

import (
	"io"
	"net"
)
// a normal closure
type regularClosure func(conn net.Conn)

// byte closure describe a closure which returns a byte array
type byteClosure func([]byte, net.Conn)

// string closure describe a closure which returns a string
type stringClosure func(string, net.Conn)

// error closure describe a closure which returns an error
type errorClosure func(error, net.Conn)

// Server represents the implementation of
// the FastSocket Protocol (Server sided)
type Server struct {
	// the framing
	frame frame
	// transfer
	transfer transfer
	// Closure for incoming text messages
	OnTextMessage stringClosure
	// Closure for incoming data messages
	OnBinaryMessage byteClosure
	// Closure for appearing errors
	OnError errorClosure
	// Closure for closed connections
	OnClose regularClosure
}
// Start starts the FastSocketServer and handles all incoming connection
func (server *Server) Start(port uint16) error {
	server.transfer = transfer{}
	server.transferClosure()
	err := server.transfer.start(port)
	if err != nil {
		return err
	}
	return nil
}
// Stop closes all tcp connections
func (server *Server) Stop() {
	err := server.transfer.stop()
	if err != nil {
		server.OnError(err, nil)
		return
	}
}
// Close is for closing a connection
func (server *Server) Close(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		server.OnError(err, conn)
		return
	}
}
// SendMessage is used to send data or string based messages to the client
func (server *Server) SendMessage(messageType messageType, data *[]byte, conn net.Conn) {
	switch messageType {
	case TextMessage:
		message, err := server.frame.create(data, TextMessage, false)
		if err != nil {
			server.OnError(err, conn)
			server.Close(conn)
			return
		}
		server.write(message, conn)
	case BinaryMessage:
		message, err := server.frame.create(data, BinaryMessage, false)
		if err != nil {
			server.OnError(err, conn)
			server.Close(conn)
			return
		}
		server.write(message, conn)
	}
}
// internal used transfer read loop
func (server *Server) transferClosure() {
	server.transfer.onConnection = func(conn net.Conn) {
		isLocked := false
		frame := frame{}
		server.frameClosures(&frame, conn)
		for {
			buffer := make([]byte, maximumLength)
			size, err := conn.Read(buffer)
			if err != nil {
				if err == io.EOF {
					server.OnClose(conn)
					return
				}
				server.OnError(err, conn)
				return
			}
			data := buffer[:size]
			switch isLocked {
			case true:
				err := frame.parse(&data)
				if err != nil {
					server.OnError(err, conn)
					server.Close(conn)
					return
				}
			case false:
				if string(data) == socketID {
					isLocked = true
					server.write(&[]byte{byte(acceptByte)}, conn)
				} else {
					server.Close(conn)
				}
			}
		}
	}
}
// internal sending function
func (server *Server) write(data *[]byte, conn net.Conn) {
	_, err := conn.Write(*data)
	if err != nil {
		server.OnError(err, conn)
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
