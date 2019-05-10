package fastsocket

import (
	"log"
	"net"
)

// byte closure describe a closure which returns a byte array
type byteClosure func([]byte, net.Conn)

// string closure describe a closure which returns a string
type stringClosure func(string, net.Conn)

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
	server.transfer.stop()
}
// SendMessage is used to send data or string based messages to the client
func (server *Server) SendMessage(messageType messageType, data *[]byte, conn net.Conn) {
	if messageType == TextMessage {
		message, err := server.frame.create(data, TextMessage, false)
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		server.write(message, conn)
	}
	if messageType == BinaryMessage {
		message, err := server.frame.create(data, BinaryMessage, false)
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
