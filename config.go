package monkey

// Config monkey 实例的相关配置
type Config struct {
	WSListernPort int // WSListernPort ws server 的监听端口
}

// NewConfig new config
func NewConfig() *Config {
	return &Config{
		WSListernPort: 10001,
	}
}

// TransportConfig transoport config
type TransportConfig struct {
	PingInterval    int // PingInterval ping的 频率 单位 Second
	WriteMaxDurtion int // WriteMaxDurtion 发送消息的最大延迟
	SendBufChanSize int // SendBufChanSize 发送缓存大小
}

// NewTransportConfig transport config
func NewTransportConfig() *TransportConfig {
	return &TransportConfig{
		PingInterval:    3,
		WriteMaxDurtion: 3,
		SendBufChanSize: 100,
	}
}
