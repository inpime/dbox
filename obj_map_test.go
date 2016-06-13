package dbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapObj_ReadWrite(t *testing.T) {
	m := NewMapObject(nil)
	mapSet(m.Map(), "a", "b")
	m.Encode()
	network := m.Bytes()

	m = NewMapObject(nil)
	m.Write(network)
	err := m.Decode()
	assert.NoError(t, err, "map decode")

	assert.Equal(t, mapString(m.Map(), "a"), "b", "not expected value")
}

func BenchmarkMapObj_ReadWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := NewMapObject(nil)
		network := m.Bytes()
		m = NewMapObject(nil)
		m.Write(network)
	}
}
