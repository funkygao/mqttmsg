package mqttmsg

const (
	// Maximum payload size in bytes (256MiB - 1B).
	MaxPayloadSize = (1 << (4 * 7)) - 1
)

const (
	RetCodeAccepted = ReturnCode(iota)
	RetCodeUnacceptableProtocolVersion
	RetCodeIdentifierRejected
	RetCodeServerUnavailable
	RetCodeBadUsernameOrPassword
	RetCodeNotAuthorized

	retCodeFirstInvalid
)

const (
	QosAtMostOnce = QosLevel(iota)
	QosAtLeastOnce
	QosExactlyOnce

	qosFirstInvalid
)

// MessageType constants.
const (
	MsgConnect = MessageType(iota + 1)
	MsgConnAck
	MsgPublish
	MsgPubAck
	MsgPubRec  // publish received (assured delivery part 1)
	MsgPubRel  // publish release (assured delivery part 2)
	MsgPubComp // publish complete (assured delivery part 3)
	MsgSubscribe
	MsgSubAck
	MsgUnsubscribe
	MsgUnsubAck
	MsgPingReq
	MsgPingResp
	MsgDisconnect

	msgTypeFirstInvalid
)
