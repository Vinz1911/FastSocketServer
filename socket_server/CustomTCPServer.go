package main

import (
	"log"
	"net"
	"strconv"
)

// Descripes a Closure which returns some
// byte array
type closure func([]byte)

// ControlCode are the Codes which
// describe the fragmented frame
type ControlCode uint8

// Opcode are Status Codes for
// the Message
type Opcode uint8

const (
	continueByte    ControlCode = 0x0
	finByte         ControlCode = 0xFF
	text            Opcode      = 0x1
	binary          Opcode      = 0x2
	connectionClose Opcode      = 0x8
	port            uint16      = 3333
	socketKey       string      = "6D8EDFD9-541C-4391-9171-AD519876B32E"
)

// FastSocket is the Custom TCP Protocol
type FastSocket struct {
	locked bool
}

func main() {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Core TCP Transfer Protocol Server Startup..., Listen on Port:", strconv.Itoa(int(port)))
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(socket net.Conn) {
	s := FastSocket{}
	f := Frame{}
	handleResponse(&f, socket)
	defer socket.Close()
	for {
		buf := make([]byte, 8192)
		size, err := socket.Read(buf)
		if err != nil {
			return
		}
		data := buf[:size]
		if s.locked == true {
			f.parse(&data)
		}
		if s.locked == false {
			if string(data) == socketKey {
				s.locked = true
				socket.Write([]byte{0xFE})
			} else {
				socket.Close()
			}
		}
	}
}

// Handle Response for Speedtest
func handleResponse(f *Frame, socket net.Conn) {
	f.onText = func(data []byte) {
		message := string(data)
		chunk := make([]byte, toInteger(message))
		frame := f.create(&chunk, binary)
		send(frame, socket)
	}
	f.onBinary = func(data []byte) {
		message := toString(len(data))
		chunk := []byte(message)
		frame := f.create(&chunk, text)
		send(frame, socket)
	}
}

func send(data *[]byte, socket net.Conn) {
	socket.Write(*data)
}

// Frame is a struct to Create and parse the
// Protocol Messages
/*
	Custom TCP Communication Protocol Framing

	+---+------------------------------+-+
	|0 1|         ... Continue         |N|
	+---+------------------------------+-+
	| O |                              |F|
	| P |         Payload Data...      |I|
	| C |                              |N|
	| O |         Payload Data...      |B|
	| D |                              |Y|
	| E |         Payload Data...      |T|
	|   |                              |E|
	|   |         Payload Data...      | |
	|   |                              | |
	+---+------------------------------+-+
*/
type Frame struct {
	cache    []byte
	onText   closure
	onBinary closure
}

func (Frame) create(data *[]byte, opcode Opcode) *[]byte {
	buffer := []byte{}
	buffer = append(buffer, byte(opcode))
	buffer = append(buffer, byte(continueByte))
	buffer = append(buffer, *data...)
	buffer = append(buffer, byte(finByte))
	return &buffer
}
func (f *Frame) parse(data *[]byte) {
	if len(*data) <= 0 {
		return
	}
	f.cache = append(f.cache, *data...)
	if (*data)[len(*data)-1] == byte(finByte) {
		if f.cache[0] == byte(text) {
			_, f.cache = f.cache[0], f.cache[1:]
			_, f.cache = f.cache[0], f.cache[1:]
			f.cache = f.cache[:len(f.cache)-1]
			f.onText(f.cache)
		}
		if f.cache[0] == byte(binary) {
			_, f.cache = f.cache[0], f.cache[1:]
			_, f.cache = f.cache[0], f.cache[1:]
			f.cache = f.cache[:len(f.cache)-1]
			f.onBinary(f.cache)
		}
		f.cache = []byte{}
	}
}

// Helper Stuff
func toInteger(msg string) int {
	value, err := strconv.Atoi(msg)
	if err != nil {
	}
	return value
}

func toString(msg int) string {
	value := strconv.Itoa(msg)
	return value
}
