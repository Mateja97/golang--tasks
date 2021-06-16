package EVM

import (
	"math"
)

type Memory struct {
	store []byte
	lastGasCost uint64
}
func NewMemory() *Memory {
	return &Memory{}
}
func (m *Memory) Set(offset uint64, val []byte) uint64{

	m.Resize(offset+uint64(len(val)))
	copy(m.store[offset:],val)
	w:=  m.CountWords()
	return w
}
func (m *Memory) GetStore() []byte {
	return m.store
}
func (m *Memory) Set8(offset uint64,byte byte) uint64 {

	m.Resize(offset)
	m.store[offset] = byte
	w:=  m.CountWords()
	return w
}

func (m* Memory) Resize(offset uint64){
	for uint64(len(m.store)) <= offset {
		m.store = append(m.store, 0)
	}

}

func (m Memory) CountWords() uint64 {
	size := uint64(len(m.store))
	if size > math.MaxUint64-31 {
		return math.MaxUint64/32 + 1
	}
	return (size + 31) / 32
}