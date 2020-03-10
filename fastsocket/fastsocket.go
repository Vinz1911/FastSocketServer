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
	OnStringMessage stringClosure
	// Closure for incoming data messages
	OnDataMessage dataClosure
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
		server.OnError(nil, err)
		return
	}
}
// Close is for closing a connection
func (server *Server) Close(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		server.OnError(conn, err)
		return
	}
	server.OnClose(conn)
}
// send string based message
func (server *Server) SendStringMessage(conn net.Conn, str string) {
	message, err := server.frame.create([]byte(str), StringMessage)
	if err != nil {
		server.OnError(conn, err)
		server.Close(conn)
		return
	}
	server.send(message, conn)
}
// send data based message
func (server *Server) SendDataMessage(conn net.Conn, data []byte) {
	message, err := server.frame.create(data, DataMessage)
	if err != nil {
		server.OnError(conn, err)
		server.Close(conn)
		return
	}
	server.send(message, conn)
}
// closures from the transfer protocol
func (server *Server) callbacks() {
	server.transfer.onMessage = func(conn net.Conn, data []byte, isLocked *bool, frame *frame) {
		switch *isLocked {
		case true:
			err := frame.parse(data, func(str string) {
				server.OnStringMessage(conn, str)
			}, func(data []byte) {
				server.OnDataMessage(conn, data)
			})
			if err != nil {
				server.OnError(conn, err)
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
		server.OnError(conn, err)
		return
	}
}