// Implementation of MQTT V3.1 encoding and decoding.
//
// See http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html
// for the MQTT protocol specification. This package does not implement the
// semantics of MQTT, but purely the encoding and decoding of its messages.
//
// Decoding Messages:
//
// Use the DecodeOneMessage function to read a Message from an io.Reader, it
// will return a Message value. The function can be implemented using the public
// API of this package if more control is required. For example:
//
//   for {
//     msg, err := mqtt.DecodeOneMessage(conn, nil)
//     if err != nil {
//       // handle err
//     }
//     switch msg := msg.(type) {
//     case *Connect:
//       // ...
//     case *Publish:
//       // ...
//       // etc.
//     }
//   }
//
// Encoding Messages:
//
// Create a message value, and use its Encode method to write it to an
// io.Writer. For example:
//
//   someData := []byte{1, 2, 3}
//   msg := &Publish{
//     Header: {
//       DupFlag: false,
//       QosLevel: QosAtLeastOnce,
//       Retain: false,
//     },
//     TopicName: "a/b",
//     MessageId: 10,
//     Payload: BytesPayload(someData),
//   }
//   if err := msg.Encode(conn); err != nil {
//     // handle err
//   }
//
// Advanced PUBLISH payload handling:
//
// The default behaviour for decoding PUBLISH payloads, and most common way to
// supply payloads for encoding, is the BytesPayload, which is a []byte
// derivative.
//
// More complex handling is possible by implementing the Payload interface,
// which can be injected into DecodeOneMessage via the `config` parameter, or
// into an outgoing Publish message via its Payload field.  Potential benefits
// of this include:
//
// * Data can be (un)marshalled directly on a connection, without an unecessary
// round-trip via bytes.Buffer.
//
// * Data can be streamed directly on readers/writers (e.g files, other
// connections, pipes) without the requirement to buffer an entire message
// payload in memory at once.
//
// The limitations of these streaming features are:
//
// * When encoding a payload, the encoded size of the payload must be known and
// declared upfront.
//
// * The payload size (and PUBLISH variable header) can be no more than 256MiB
// minus 1 byte. This is a specified limitation of MQTT v3.1 itself.
package mqttmsg
