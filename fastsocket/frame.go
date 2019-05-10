package fastsocket

import (
	"encoding/binary"
	"errors"
)

// Frame is a struct to Create and parse the
// Protocol Messages
// Custom TCP Communication Protocol Framing
//
// 0                   1       N
// +-------------------+-------+
// |0|1|2 3 4 5 6 7 8 9|0 1 2 3|
// +-+-+---------------+-------+
// |F|O| Payload length|PAYLOAD|
// |I|P|      (8)      |  (N)  |
// |N|C|               |       |
// +-+-+---------------+-------+
// : Payload Data continued ...:
// + - - - - - - - - - - - - - +
// | Payload Data continued ...|
// +---------------------------+
//
type frame struct {
	cache    []byte
	onText   func([]byte)
	onBinary func([]byte)
}
// Create a FastSocket Protocol compliant frame
// this functions add the necessary bytes to the buffer
func (*frame) create(data *[]byte, messageType messageType, isFin bool) (*[]byte, error) {
	var buffer []byte
	var size = make([]byte, 8)
	binary.LittleEndian.PutUint64(size, uint64(len(*data) + overheadSize))
	if isFin {
		buffer = append(buffer, byte(finByte))
	} else {
		buffer = append(buffer, byte(continueByte))
	}
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
		return errors.New("zero data")
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
		switch slice[1] {
		case byte(TextMessage):
			message, err := frame.trimmedFrame(&slice)
			if err != nil {
				return err
			}
			frame.onText(message)
		case byte(BinaryMessage):
			message, err := frame.trimmedFrame(&slice)
			if err != nil {
				return err
			}
			frame.onBinary(message)
		default:
			return errors.New("invalid opcode")
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
	size := frame.cache
	size = size[2:10]
	return int(binary.LittleEndian.Uint64(size))
}

// helper function to trim a frame into a message
func (*frame) trimmedFrame(data *[]byte) ([]byte, error) {
	if len(*data) < overheadSize {
	return nil, errors.New("cannot trim frame")
	}
	var trimmed = *data
	if len(*data) >= 10 {
		trimmed = trimmed[10:]
	} else {
		return nil, errors.New("cannot trim frame")
	}
	return trimmed, nil
}