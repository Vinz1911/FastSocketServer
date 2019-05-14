package fastsocket

import (
	"crypto/tls"
	"net"
	"strconv"
)

// transfer is a raw tcp socket
type transfer struct {
	// the tcp listener
	listener net.Listener
	// go routine closure with the connection
	onConnection func(net.Conn)
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
		if err != nil {
			return  err
		}
	case TLSTransfer:
		cer, err := tls.LoadX509KeyPair(transfer.certPath, transfer.keyPath)
		if err != nil {
			return err
		}
		config := &tls.Config{Certificates: []tls.Certificate{cer}}
		transfer.listener, err = tls.Listen("tcp", ":"+strconv.Itoa(int(port)), config)
		if err != nil {
			return err
		}
	}
	defer transfer.listener.Close()
	for {
		connection, err := transfer.listener.Accept()
		if err != nil {
			return err
		}
		go transfer.onConnection(connection)
	}
}
// invalidate all current running tcp connections
func (transfer *transfer) stop() error {
	err := transfer.listener.Close()
	if err != nil {
		return err
	}
	return nil
}