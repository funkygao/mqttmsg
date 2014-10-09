package mqtt

import (
	"bytes"
	"io"
)

// Header contains the common attributes of all messages. Some attributes are
// not applicable to some message types.
// MessageType and remaining length is outside of Header
type Header struct {
	DupFlag  bool
	Retain   bool
	QosLevel QosLevel
}

func (hdr *Header) Encode(w io.Writer, msgType MessageType, remainingLength int32) error {
	buf := new(bytes.Buffer)
	err := hdr.encodeInto(buf, msgType, remainingLength)
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	return err
}

func (hdr *Header) encodeInto(buf *bytes.Buffer, msgType MessageType, remainingLength int32) error {
	if !hdr.QosLevel.IsValid() {
		return badQosError
	}
	if !msgType.IsValid() {
		return badMsgTypeError
	}

	val := byte(msgType) << 4
	val |= (boolToByte(hdr.DupFlag) << 3)
	val |= byte(hdr.QosLevel) << 1
	val |= boolToByte(hdr.Retain)
	buf.WriteByte(val)
	encodeLength(remainingLength, buf)
	return nil
}

func (hdr *Header) Decode(r io.Reader) (msgType MessageType, remainingLength int32, err error) {
	defer func() {
		err = recoverError(err, recover())
	}()

	var buf [1]byte

	if _, err = io.ReadFull(r, buf[:]); err != nil {
		return
	}

	byte1 := buf[0]
	msgType = MessageType(byte1 & 0xF0 >> 4)

	*hdr = Header{
		DupFlag:  byte1&0x08 > 0,
		QosLevel: QosLevel(byte1 & 0x06 >> 1),
		Retain:   byte1&0x01 > 0,
	}

	remainingLength = decodeLength(r)

	return
}

type QosLevel uint8

func (qos QosLevel) IsValid() bool {
	return qos < qosFirstInvalid
}

func (qos QosLevel) HasId() bool {
	return qos == QosAtLeastOnce || qos == QosExactlyOnce
}

type ReturnCode uint8

func (rc ReturnCode) IsValid() bool {
	return rc >= RetCodeAccepted && rc < retCodeFirstInvalid
}

// DecoderConfig provides configuration for decoding messages.
type DecoderConfig interface {
	// MakePayload returns a Payload for the given Publish message. r is a Reader
	// that will read the payload data, and n is the number of bytes in the
	// payload. The Payload.ReadPayload method is called on the returned payload
	// by the decoding process.
	MakePayload(msg *Publish, r io.Reader, n int) (Payload, error)
}

type DefaultDecoderConfig struct{}

func (c DefaultDecoderConfig) MakePayload(msg *Publish, r io.Reader, n int) (Payload, error) {
	return make(BytesPayload, n), nil
}

// ValueConfig always returns the given Payload when MakePayload is called.
type ValueConfig struct {
	Payload Payload
}

func (c *ValueConfig) MakePayload(msg *Publish, r io.Reader, n int) (Payload, error) {
	return c.Payload, nil
}

// DecodeOneMessage decodes one message from r. config provides specifics on
// how to decode messages, nil indicates that the DefaultDecoderConfig should
// be used.
func DecodeOneMessage(r io.Reader, config DecoderConfig) (msg Message, err error) {
	var hdr Header
	var msgType MessageType
	var packetRemaining int32
	msgType, packetRemaining, err = hdr.Decode(r)
	if err != nil {
		return
	}

	msg, err = NewMessage(msgType)
	if err != nil {
		return
	}

	if config == nil {
		config = DefaultDecoderConfig{}
	}

	return msg, msg.Decode(r, hdr, packetRemaining, config)
}

// NewMessage creates an instance of a Message value for the given message
// type. An error is returned if msgType is invalid.
func NewMessage(msgType MessageType) (msg Message, err error) {
	switch msgType {
	case MsgConnect:
		msg = new(Connect)
	case MsgConnAck:
		msg = new(ConnAck)
	case MsgPublish:
		msg = new(Publish)
	case MsgPubAck:
		msg = new(PubAck)
	case MsgPubRec:
		msg = new(PubRec)
	case MsgPubRel:
		msg = new(PubRel)
	case MsgPubComp:
		msg = new(PubComp)
	case MsgSubscribe:
		msg = new(Subscribe)
	case MsgUnsubAck:
		msg = new(UnsubAck)
	case MsgSubAck:
		msg = new(SubAck)
	case MsgUnsubscribe:
		msg = new(Unsubscribe)
	case MsgPingReq:
		msg = new(PingReq)
	case MsgPingResp:
		msg = new(PingResp)
	case MsgDisconnect:
		msg = new(Disconnect)
	default:
		return nil, badMsgTypeError
	}

	return
}
