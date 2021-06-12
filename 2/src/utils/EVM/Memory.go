package EVM

type Memory struct {
	store []byte
	lastGasCost uint64
}
func NewMemory() *Memory {
	return &Memory{store: make([]byte,256)}
}
func (m *Memory) Set(offset uint64, val []byte) {

	if offset+uint64(len(val)) > uint64(len(m.store)) {
		panic("invalid memory: store empty")
	}
	copy(m.store[offset:],val)
}
func (m *Memory) GetStore() []byte {
	return m.store
}
func (m *Memory) SetStore(offset uint64,byte byte) {
	m.store[offset] = byte
}

