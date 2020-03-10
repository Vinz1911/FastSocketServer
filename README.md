<div align="center">
    <h1>
        <br>
            <a href="https://github.com/Vinz1911/FastSocketServer"><img src="https://github.com/Vinz1911/FastSocketServer/blob/develop/.fastsocketserver.svg" alt="FastSocket" width="500"></a>
        <br>
        <br>
            FastSocket Server
        <br>
    </h1>
</div>

`FastSocket` is a proprietary bi-directional message based communication protocol on top of TCP (optionally over other layers in the future). This is the server-sided implementation of the FastSocket Protocol. It's optimized for maximum speed and performance. Socket is the FastSocket Protocol backend implementation with a limit feature set in starting the server at a specific port and sending and receiving messages.
 
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
    "github.com/vinz1911/socket"
)

// get fastsocket server
func main() {
    server := fastsocket.Server{}
}
```

## Closures:
```go
server.OnBinaryMessage = func(data []byte, socket net.Conn) {
    // called if binary message was received
}
server.OnTextMessage = func(str string, socket net.Conn) {
    // called if text message was received
}
```

## Send Messages:
```go
// send a text or binary message to the client
// you need to transform a string into []byte before sending
server.SendMessage(messageType messageType, data []byte, socket net.Conn)
```

## Start Server:
```go
// start the server at a specific tcp port
server.Start(uint16)
```
