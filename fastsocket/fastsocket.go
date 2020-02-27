// Copyright 2019 Vinzenz Weist. All rights reserved.
// Use of this source code is risked by yourself.
// license that can be found in the LICENSE file.

// FastSocket is a proprietary communication protocol directly
// written on top of TCP. It's a message based protocol which allows you
// to send text and binary based messages. The protocol is so small it have
// only 5 Bytes overhead per message, the handshake is done directly on TCP level.
// The motivation behind this protocol was, to use it as `Speedtest Protocol`, a
// low level TCP communication protocol to measure TCP throughput performance. -> FastSockets is the answer
// FastSocket allows to enter all possible TCP Options if needed and is completely non-blocking and async,
// thanks to golang's go routine
package fastsocket

import (
	"net"
)
// Server represents the implementation of
// the FastSocket Protocol (Server sided)
type Server struct {
	// the framing
	frame frame
	// transfer
	transfer transfer
	// Closure if new connection comes in
	OnReady regularClosure
	// Closure for incoming text messages
	OnTextMessage stringClosure
	// Closure for incoming data messages
	OnBinaryMessage byteClosure
	// Closure for appearing errors
	OnError errorClosure
	// Closure for closed connections
	OnClose regularClosure
	// the path to the cert file for tls
	CertPath string
	// the path to the key file for tls
	KeyPath string
}
// Start starts the FastSocketServer and handles all incoming connection
func (server *Server) Start(transferType transferType, port uint16) error {
	server.transfer = transfer{}
	server.transfer.certPath = server.CertPath
	server.transfer.keyPath = server.KeyPath
	server.callbacks()
	err := server.transfer.start(transferType, port)
	if err != nil { return err }
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
	server.OnClose(conn)
}
// SendMessage is used to send data or string based messages to the client
func (server *Server) SendMessage(messageType messageType, data []byte, conn net.Conn) {
	switch messageType {
	case TextMessage:
		message, err := server.frame.create(data, TextMessage)
		if err != nil {
			server.OnError(err, conn)
			server.Close(conn)
			return
		}
		server.send(message, conn)
	case BinaryMessage:
		message, err := server.frame.create(data, BinaryMessage)
		if err != nil {
			server.OnError(err, conn)
			server.Close(conn)
			return
		}
		server.send(message, conn)
	}
}
// closures from the transfer protocol
func (server *Server) callbacks() {
	server.transfer.onMessage = func(data []byte, conn net.Conn, isLocked *bool, frame *frame) {
		switch *isLocked {
		case true:
			err := frame.parse(data, func(str string) {
				server.OnTextMessage(str, conn)
			}, func(data []byte) {
				server.OnBinaryMessage(data, conn)
			})
			if err != nil {
				server.OnError(err, conn)
				server.Close(conn)
				return
			}
		case false:
			if isUUID(string(data)) {
				sha256 := generateSHA256(string(data))
				mapped := sha256[:]
				server.send(mapped, conn)
				*isLocked = true
			} else {
				server.Close(conn)
			}
		}
	}
	server.transfer.onReady = server.OnReady
	server.transfer.onClose = server.OnClose
	server.transfer.onError = server.OnError
}
// internal sending function
func (server *Server) send(data []byte, conn net.Conn) {
	_, err := conn.Write(data)
	if err != nil {
		server.OnError(err, conn)
		return
	}
}