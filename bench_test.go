package mqttmsg

import (
	"io/ioutil"
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	b.ReportAllocs()
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	msg := Publish{
		Header: Header{
			DupFlag:  false,
			QosLevel: QosAtLeastOnce,
			Retain:   false,
		},
		TopicName: "a/b",
		MessageId: 10,
		Payload:   BytesPayload(data),
	}
	for i := 0; i < b.N; i++ {
		msg.Encode(ioutil.Discard)
	}
	b.SetBytes(int64(len(data)))
}
