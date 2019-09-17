package fastsocket

import (
	"encoding/binary"
	"errors"
)

// 0                   1       N
// +-------------------+-------+
// |0|1|2 3 4 5 6 7 8 9|0 1 2 3|
// +-+-+---------------+-------+
// |F|O| FRAME LENGTH  |PAYLOAD|
// |I|P|     (8)       |  (N)  |
// |N|C|               |       |
// +-+-+---------------+-------+
// : Payload Data continued ...:
// + - - - - - - - - - - - - - +
// | Payload Data continued ...|
// +---------------------------+
//
// This describes the framing protocol.
// - FIN: 0x3
//      - The first byte is used to inform the the other side, that the
//        connection is finished and can be closed, this is used to prevent
//        that a connection will be closed but there are unread bytes on the connection
// - OPC:
//      - 0x0: this is the continue byte (currently a placeholder)
//      - 0x1: this is the string byte which is used for string based messages
//      - 0x2: this is the data byte which is used for data based messages
//      - 0x3: this is the fin byte, which is part of OPC but is on the first place in the protocol
//      - 0x6: this is the accept byte and is used by the handshake
//      - 0x7 - 0xF: this bytes are reserved
// - FRAME LENGTH:
//      - this uses 8 bytes to store the entire frame size as a big endian uint64 value
// - PAYLOAD:
//      - continued payload data

type frame struct {
	cache    []byte
	onText   func(string)
	onBinary func([]byte)
}
// Create a FastSocket Protocol compliant frame
// this functions add the necessary bytes to the buffer
func (*frame) create(data *[]byte, messageType messageType) (*[]byte, error) {
	var buffer []byte
	var size = make([]byte, 8)
	binary.BigEndian.PutUint64(size, uint64(len(*data) + overheadSize))
	buffer = append(buffer, byte(messageType))
	buffer = append(buffer, size...)
	buffer = append(buffer, *data...)
	if len(buffer) > maximumContentLength {
		return nil, errors.New("write buffer overflow")
	}
	return &buffer, nil
}
// Parse received Data into a FastSocket compliant
// frame/message
func (frame *frame) parse(data *[]byte) error {
	if len(*data) <= 0 {
		return errors.New("zero data fault")
	}
	frame.cache = append(frame.cache, *data...)
	if len(frame.cache) > maximumContentLength {
		return errors.New("read buffer overflow")
	}
	if len(frame.cache) < overheadSize {
		return nil
	}
	if len(frame.cache) < frame.contentSize() {
		return nil
	}
	for len(frame.cache) >= frame.contentSize() && frame.contentSize() != 0 {
		slice := frame.cache[:frame.contentSize()]
		switch slice[0] {
		case byte(TextMessage):
			message, err := frame.trimFrame(slice)
			trimmed := string(message)
			if err != nil {
				return err
			}
			frame.onText(trimmed)
		case byte(BinaryMessage):
			trimmed, err := frame.trimFrame(slice)
			if err != nil {
				return err
			}
			frame.onBinary(trimmed)
		default:
			return errors.New("invalid operational code")
		}
		if len(frame.cache) > frame.contentSize() {
			frame.cache = frame.cache[frame.contentSize():]
		} else {
			frame.cache = []byte{}
		}
	}
	return nil
}

// helper function to parse the content size of a frame
func (frame *frame) contentSize() int {
	if len(frame.cache) < overheadSize {
		return 0
	}
	size := frame.cache[1:overheadSize]
	return int(binary.BigEndian.Uint64(size))
}

// helper function to trim a frame into a message
func (*frame) trimFrame(frame []byte) ([]byte, error) {
	if len(frame) < overheadSize {
		return nil, errors.New("cannot trim frame")
	}
	data := frame[overheadSize:]
	return data, nil
}