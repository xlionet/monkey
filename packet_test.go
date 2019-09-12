package monkey

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GatePacketer_Pack(t *testing.T) {
	maybeMsg := struct {
		Name string
		Age  int32
		Msg  []byte
	}{
		Name: "yanxi",
		Age:  20,
		Msg:  []byte("yanxi is young"),
	}
	m, _ := json.Marshal(maybeMsg)
	packet := JSONGatePacket{PacketID: 1, Data: m}

	bp, _ := json.Marshal(packet)
	var base JSONGatePacket
	err := JSONPacketer.Unpack(bp, &base)
	assert.Nil(t, err)

	t.Logf("base pack is : %#v", base)

}
