package monkey

import "encoding/json"

// GatePacketer gate 数据的打包器
type GatePacketer interface {
	Pack(int, int, interface{}) (*Envelope, error)
	Unpack([]byte, interface{}) error
}

// JSONGatePacket 默认网关解包器
type JSONGatePacket struct {
	PacketID int             `json:"packet_id"`
	Data     json.RawMessage `json:"data"`
}

// Pack 打包gate 包数据
func (p *JSONGatePacket) Pack(msgType, packetID int, data interface{}) (*Envelope, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var pack = &JSONGatePacket{
		PacketID: packetID,
		Data:     d,
	}
	msg, err := json.Marshal(pack)
	if err != nil {
		return nil, err
	}

	return &Envelope{T: msgType, Msg: msg}, nil
}

// Unpack 解析
func (p *JSONGatePacket) Unpack(raw []byte, v interface{}) error {
	return json.Unmarshal(raw, v)
}

// JSONPacketer 默认的json 解析器
var JSONPacketer = JSONGatePacket{}
