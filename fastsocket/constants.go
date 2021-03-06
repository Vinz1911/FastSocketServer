// Copyright 2019 Vinzenz Weist. All rights reserved.
// Use of this source code is risked by yourself.
// license that can be found in the LICENSE file.
package fastsocket

// operational codes are used to
// control the framing, handles handshake and more
type operationalCode uint8

// the message type for framing
type messageType uint8

// transfer type
type transferType int

const (
	// tcp transfer type
	TCPTransfer transferType = 0
	// tls transfer type
	TLSTransfer transferType = 1
	// continueByte is a placeholder `UNUSED`
	continueByte operationalCode = 0x0
	// text holds the byte for a text message
	StringMessage messageType = 0x1
	// binary holds the byte for a binary message
	DataMessage messageType = 0x2
	// finByte holds the ControlCode for `end of a message`
	finByte operationalCode = 0x03
	// maximumLength is the maximum buffer read length
	maximumLength int = 8192
	// maximum frame size
	maximumFrameLength int = 16_777_216
	// overhead
	overheadSize int = 5
)