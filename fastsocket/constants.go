package fastsocket


// operational codes are used to
// control the framing, handles handshake and more
type operationalCode uint8

// the message type for framing
type messageType uint8

const (
	// continueByte is a placeholder `UNUSED`
	 continueByte operationalCode = 0x0
	// text holds the byte for a text message
	TextMessage messageType = 0x1

	// binary holds the byte for a binary message
	BinaryMessage messageType = 0x2

	// finByte holds the ControlCode for `end of a message`
	finByte operationalCode = 0x03

	// acceptByte is for ACK if the handshake succeed
	acceptByte operationalCode = 0x06

	// connectionClose holds the byte for convenient closed connections `UNUSED`
	connectionClose operationalCode = 0x8

	// socketID is the unique ID to identify protocol conformance connections (used for the handshake)
	socketID string = "6D8EDFD9-541C-4391-9171-AD519876B32E"

	// maximumLength is the maximum buffer read length
	maximumLength int = 8192

	// maximum frame size
	maximumContentLength int = 16777216

	// overhead
	overheadSize int = 10
)