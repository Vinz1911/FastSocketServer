package fastsocket

import "net"
// a normal closure
type regularClosure func(conn net.Conn)

// byte closure describe a closure which returns a byte array
type byteClosure func([]byte, net.Conn)

// string closure describe a closure which returns a string
type stringClosure func(string, net.Conn)

// error closure describe a closure which returns an error
type errorClosure func(error, net.Conn)

// transfer closure for the transfer
type transferClosure func([]byte, net.Conn, *bool, *frame)