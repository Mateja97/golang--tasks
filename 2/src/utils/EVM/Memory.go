package EVM

import (
	"math"
)

type Memory struct {
	store       []byte
	lastGasCost uint64
}

func NewMemory() *Memory {
	return &Memory{}
}
func (m *Memory) Set(offset uint64, val []byte) uint64 {

	m.Resize(offset) //resize memory if needed
	copy(m.store[offset:], val)
	w := m.CountWords() //count words in new memory
	return w
}
func (m *Memory) Set8(offset uint64, byte byte) uint64 {

	m.Resize(offset) //resize memory if needed
	m.store[offset] = byte
	w := m.CountWords() //count words in new memory
	return w
}

func (m *Memory) Resize(reqSize uint64) {
	for uint64(len(m.store)) <= reqSize { // if offset is bigger then current length - resize until data could be stored with new offset
		m.store = append(m.store, 0)
	}

}

func (m Memory) CountWords() uint64 {
	size := uint64(len(m.store))
	if size > math.MaxUint64-31 {
		return math.MaxUint64/32 + 1
	}
	return (size + 31) / 32 // Calculate number of words (+31 - to count words up to 32 bytes as 1 word)
}
