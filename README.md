<div align="center">
    <h1>
        <br>
            <a href="https://github.com/Vinz1911/FastSocketServer"><img src="https://github.com/Vinz1911/FastSocketServer/blob/master/.fastsocketserver.svg" alt="FastSocket" width="500"></a>
        <br>
        <br>
            FastSocket Server
        <br>
    </h1>
</div>

`FastSocket` is a proprietary bi-directional message based communication protocol on top of TCP (optionally over other layers in the future). The idea behind this project was, to create a TCP communication like the [WebSocket Protocol](https://tools.ietf.org/html/rfc6455) with less overhead and the ability to track every 8192 bytes read or written on the socket without waiting for the whole message to be transmitted. This allows it to use it as **protocol for speed tests** for measuring the TCP throughput performance. Our server-sided implementation is written in [golang](https://golang.org/) and it's optimized for maximum speed and performance.

The client sided (Swift) implementation of the FastSocket Protocol can be found here: [FastSocket](https://github.com/Vinz1911/FastSocket).

## Features:
- [X] send and receive text and data messages
- [X] async, non-blocking & very fast
- [X] go routines handle every single tcp connection
- [X] zer0 dependencies, native go implementation
- [X] maximum frame size 16777216 bytes (with overhead)
- [X] content length base framing instead of fin byte termination
- [X] send/receive multiple messages at once
- [X] TLS support

## License:
[![License](https://img.shields.io/badge/license-GPLv3-blue.svg?longCache=true&style=flat)](https://github.com/Vinz1911/FastSocketServer/blob/master/LICENSE)

## Golang Version:
[![Golang 1.14](https://img.shields.io/badge/Golang-1.14-00ADD8.svg?logo=go&style=flat)](https://golang.org) [![Golang 1.14](https://img.shields.io/badge/Packages-Support-00ADD8.svg?logo=go&style=flat)](https://golang.org)

## Install:
```shell script
go get github.com/vinz1911/fastsocketserver
```

## Import:
```go
package main

// the net package must also be imported to
// send messages back on a specific socket
import (
    "net"
    "github.com/vinz1911/fastsocketserver"
)

// get fastsocket server
func main() {
    server := fastsocket.Server{}
}
```

## Closures:
```go
server.OnDataMessage = func(CONN net.Conn, DATA []byte) {
    // called when a binary message was received
}
server.OnStringMessage = func(CONN net.Conn, STRING string) {
    // called when a text message was received
}
```

## Send Messages:
```go
// send a binary message to the client
// CONN: net.Conn Object
// MESSAGE: binary based message
server.SendDataMessage(CONN, MESSAGE)

// send a text message to the client
// CONN: net.Conn Object
// MESSAGE: text based message
server.SendStringMessage(CONN, MESSAGE)
```

## Start Server:
```go
// start the server at a specific tcp port with a plain TCP or TLS connection
// PORT: port of host (Uint16 value)
// TRANSFER_TYPES:
//  - fastsocket.TCPTransfer (for unsecure connection)
//  - fastsocket.TLSTransfer (for secure connection, requires load of ssl certs)
server.Start(TRANSFER_TYPE, PORT)
```

## Author:
👨🏼‍💻 [Vinzenz Weist](https://github.com/Vinz1911)

This is my heart project, it's made with a lot of love and dedication ❤️