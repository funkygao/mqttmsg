package mqttmsg

import (
	"bytes"
	"github.com/funkygao/golib/recycler"
)

var (
	bufferPoolGet, bufferPoolGive = recycler.New(100, func() interface{} {
		return new(bytes.Buffer)
	})
	bytesPoolGet, bytesPoolPut = recycler.New(100, func() interface{} {
		return make([]byte, 100<<10) // max 100KB
	})
)
