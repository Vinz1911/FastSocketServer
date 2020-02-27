package fastsocket

import "net"
// a normal closure
type regularClosure func(conn net.Conn)

// byte closure describe a closure which returns a byte array
type dataClosure func(net.Conn, []byte)

// string closure describe a closure which returns a string
type stringClosure func(net.Conn, string)

// error closure describe a closure which returns an error
type errorClosure func(net.Conn, error)

// transfer closure for the transfer
type transferClosure func(net.Conn, []byte, *bool, *frame)