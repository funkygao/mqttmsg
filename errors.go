package mqttmsg

import (
	"errors"
)

var (
	badMsgTypeError        = errors.New("Invalid message type")
	badQosError            = errors.New("Invalid QoS level")
	badWillQosError        = errors.New("Invalid will QoS")
	badLengthEncodingError = errors.New("remaining length field exceeded maximum of 4 bytes")
	badReturnCodeError     = errors.New("Invalid return code")
	dataExceedsPacketError = errors.New("data exceeds packet length")
	msgTooLongError        = errors.New("message is too long")
)

// ConnectionErrors is an array of errors corresponding to the
// Connect return codes specified in the specification.
var ConnectionErrors = [6]error{
    nil, // RetCodeAccepted
    errors.New("Connection Refused: unacceptable protocol version"), // RetCodeUnacceptableProtocolVersion
    errors.New("Connection Refused: identifier rejected"),           // RetCodeIdentifierRejected
    errors.New("Connection Refused: server unavailable"),            // RetCodeServerUnavailable
    errors.New("Connection Refused: bad user name or password"),     // RetCodeBadUsernameOrPassword
    errors.New("Connection Refused: not authorized"),                // RetCodeNotAuthorized
}

// panicErr wraps an error that caused a problem that needs to bail out of the
// API, such that errors can be recovered and returned as errors from the
// public API.
type panicErr struct {
	err error
}

func (p panicErr) Error() string {
	return p.err.Error()
}

func raiseError(err error) {
	panic(panicErr{err})
}

// recoverError recovers any panic in flight and, iff it's an error from
// raiseError, will return the error. Otherwise re-raises the panic value.
// If no panic is in flight, it returns existingErr.
//
// This must be used in combination with a defer in all public API entry
// points where raiseError could be called.
func recoverError(existingErr error, recovered interface{}) error {
	if recovered != nil {
		if pErr, ok := recovered.(panicErr); ok {
			return pErr.err
		} else {
			panic(recovered)
		}
	}
	return existingErr
}
