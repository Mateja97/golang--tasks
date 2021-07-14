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

	m.Resize(offset + 32) //resize memory if needed
	l := uint64(len(val))
	if l < 32 {
		newVal := make([]byte, 32)
		copy(newVal[32-l:], val)
		copy(m.store[offset:], newVal)

	} else {
		copy(m.store[offset:], val)
	}
	w := CountWords(m.store) //count words in new memory

	return w
}
func (m *Memory) Set8(offset uint64, val byte) uint64 {

	m.Resize(offset + 1) //resize memory if needed
	m.store[offset] = val
	w := CountWords(m.store) //count words in new memory
	return w
}

func (m *Memory) Resize(reqSize uint64) {

	if uint64(len(m.store)) < reqSize { // if offset is bigger then current length - resize data to be stored with new offset
		newMemory := make([]byte, reqSize-uint64(len(m.store)))
		m.store = append(m.store, newMemory...)
	}
}

func CountWords(val []byte) uint64 {
	size := uint64(len(val))
	if size > math.MaxUint64-31 {
		return math.MaxUint64/32 + 1
	}
	return (size + 31) / 32 // Calculate number of words (+31 - to count words up to 32 bytes as 1 word)
}
