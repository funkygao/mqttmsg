package mqttmsg

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func BenchmarkEncodePublish(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`{"payload":{"330":{"uid":53,"march_id":330,"city_id":53,"opp_uid":0,"world_id":1,"type":"encamp","start_x":72,"start_y":64,"end_x":80,"end_y":78,"start_time":1412999095,"end_time":1412999111,"speed":1,"state":"marching","alliance_id":0}}`)
	msg := Publish{
		Header: Header{
			DupFlag:  false,
			QosLevel: QosAtLeastOnce,
			Retain:   false,
		},
		TopicName: "user/134568765",
		MessageId: 10,
		Payload:   BytesPayload(data),
	}
	for i := 0; i < b.N; i++ {
		msg.Encode(ioutil.Discard)
	}
	b.SetBytes(int64(len(data)))
}

func BenchmarkDecodeConnect(b *testing.B) {
	b.ReportAllocs()

	m := &Connect{
		ProtocolName:    "MQIsdp",
		ProtocolVersion: 3,
		UsernameFlag:    true,
		PasswordFlag:    true,
		WillRetain:      false,
		WillQos:         1,
		WillFlag:        true,
		CleanSession:    true,
		KeepAliveTimer:  10,
		ClientId:        "xixihaha",
		WillTopic:       "topic",
		WillMessage:     "message",
		Username:        "name",
		Password:        "pwd",
	}

	buf := new(bytes.Buffer)
	m.Encode(buf)
	for i := 0; i < b.N; i++ {
		DecodeOneMessage(buf, nil)
	}
	b.SetBytes(int64(m.Bytes()))
}

func BenchmarkDecodePublish(b *testing.B) {
	b.ReportAllocs()

	data := []byte(`{"payload":{"330":{"uid":53,"march_id":330,"city_id":53,"opp_uid":0,"world_id":1,"type":"encamp","start_x":72,"start_y":64,"end_x":80,"end_y":78,"start_time":1412999095,"end_time":1412999111,"speed":1,"state":"marching","alliance_id":0}}`)
	m := Publish{
		Header: Header{
			DupFlag:  false,
			QosLevel: QosAtLeastOnce,
			Retain:   false,
		},
		TopicName: "user/134568765",
		MessageId: 10,
		Payload:   BytesPayload(data),
	}

	buf := new(bytes.Buffer)
	m.Encode(buf)
	for i := 0; i < b.N; i++ {
		DecodeOneMessage(buf, nil)
	}
	b.SetBytes(int64(m.Bytes()))
}
