package net

// Cmd returns cmd
func (m *Message) Cmd() string {
	return m.cmd
}

// Length returns length
func (m *Message) Length() uint32 {
	return m.length
}

// Payload returns payload
func (m *Message) Payload() []byte {
	return m.payload
}
