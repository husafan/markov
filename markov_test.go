package markov_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	. "github.com/husafan/markov"
	"github.com/stretchr/testify/assert"
)

func TestUint16State(t *testing.T) {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, uint16(32))

	state := Uint16State(32)
	assert.Equal(t, "32", state.Value())
	assert.Equal(t, uint64(2), state.Size())
	assert.Equal(t, bytes, state.Bytes())
}

func TestStringState(t *testing.T) {
	hello := bytes.NewBufferString("hello").Bytes()

	state := StringState("hello")
	assert.Equal(t, "hello", state.Value())
	assert.Equal(t, uint64(5), state.Size())
	assert.Equal(t, hello, state.Bytes())
}

func TestNormalizingRowReturnsNotOkForNonValue(t *testing.T) {
	row := NewNormalizingRow()
	state := Uint16State(32)
	value, ok := row.StateWeight(state)
	assert.Equal(t, uint64(0), value)
	assert.Equal(t, false, ok)
}

func TestNormalizingRowSize(t *testing.T) {
	row := NewNormalizingRow()
	assert.Equal(t, uint64(4), row.Size())

	// A Uint16State has a size of 2 bytes. So after adding it to the row
	// the new size should be 6.
	uintState := Uint16State(32)
	assert.Equal(t, uint64(6), row.AddState(uintState))

	// The following StringState has 1 byte per character, so the new total
	// size should be 6 + 5 = 11.
	stringState := StringState("Hello")
	assert.Equal(t, uint64(11), row.AddState(stringState))
}
