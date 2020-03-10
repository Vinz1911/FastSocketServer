// Copyright 2019 Vinzenz Weist. All rights reserved.
// Use of this source code is risked by yourself.
// license that can be found in the LICENSE file.
package fastsocket

import (
	"crypto/tls"
	"io"
	"net"
	"strconv"
)
// transfer is a raw tcp socket
type transfer struct {
	// the tcp listener
	listener net.Listener
	// closure if connection is ready
	onReady regularClosure
	// closure to provide data
	onMessage transferClosure
	// closure if connection is closed
	onClose regularClosure
	// closure if connection is errored
	onError errorClosure
	// path for tls cert
	certPath string
	// path for tls key
	keyPath string
}
// start the network listener on specific port
func (transfer *transfer) start(transferType transferType, port uint16) error {
	var err error
	switch transferType {
	case TCPTransfer:
		transfer.listener, err = net.Listen("tcp", ":"+strconv.Itoa(int(port)))
		if err != nil { return  err }
	case TLSTransfer:
		cer, err := tls.LoadX509KeyPair(transfer.certPath, transfer.keyPath)
		if err != nil { return err }
		config := &tls.Config{Certificates: []tls.Certificate{cer}}
		transfer.listener, err = tls.Listen("tcp", ":"+strconv.Itoa(int(port)), config)
		if err != nil { return err }
	}
	defer transfer.listener.Close()
	for {
		connection, err := transfer.listener.Accept()
		if err != nil { return err }
		go transfer.readLoop(connection)
	}
}
// read the data from the socket
func (transfer *transfer) readLoop(conn net.Conn) {
	isLocked := false
	frame := frame{}
	transfer.onReady(conn)
	buffer := make([]byte, maximumLength)
	for {
		size, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				transfer.onClose(conn)
				return
			}
			transfer.onError(conn, err)
			return
		}
		data := buffer[:size]
		transfer.onMessage(conn, data, &isLocked, &frame)
	}
}
// invalidate all current running tcp connections
func (transfer *transfer) stop() error {
	err := transfer.listener.Close()
	if err != nil { return err }
	return nil
}