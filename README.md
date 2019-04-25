# FastSocketServer

## About:
FastSocketServer is a golang package which implements the `FastSocket Protocol`.
The FastSocket Protocol is a message based tcp communication protocol, it is probably the 
fastest message based communication protocol with the lowest overhead.
 
## Install:
```
go get github.com/vinz1911/fastsocketserver
```

## Usage:
```go
package main

import (
"net"

"github.com/vinz1911/fastsocketserver"
)

func main() {
server := fastsocket.Server{}
server.Start(8081)
}
```

## Example:
Go into the example folder to see a example implementation of the FastSocket Server.
I used it to create a `Speedtest Backend`
