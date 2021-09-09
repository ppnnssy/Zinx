package znet

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

//提供一个创建方法，工厂模式
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		Data:    data,
		DataLen: uint32(len(data)),
	}
}

// GetMsgId 获取消息ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// GetMsgLen 获取消息长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// GetMsgData 获取消息内容
func (m *Message) GetMsgData() []byte {
	return m.Data
}

// SetMsgId 设置消息ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

// SetMsgLen 设置消息长度
func (m *Message) SetMsgLen(len uint32) {
	m.DataLen = len
}

// SetMsgData 设置消息内容
func (m *Message) SetMsgData(data []byte) {
	m.Data = data
}
