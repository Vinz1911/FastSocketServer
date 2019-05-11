package fastsocket

import (
	"net"
	"strconv"
)

// transfer is a raw tcp socket
type transfer struct {
	// the tcp listener
	listener net.Listener
	// go routine closure with the connection
	onConnection func(net.Conn)
}
// start the network listener
func (transfer *transfer) start(port uint16) error {
	var err error
	transfer.listener, err = net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		return  err
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