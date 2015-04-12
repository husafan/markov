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
	value := row.StateWeight(state)
	assert.Equal(t, float64(0), value)
}

func TestNormalizingRowSize(t *testing.T) {
	row := NewNormalizingRow()
	assert.Equal(t, uint64(4), row.Size())

	// A Uint16State has a size of 2 bytes plus a 4 byte counts. So after
	// adding it to the row the new size should be 10.
	uintState := Uint16State(32)
	assert.Equal(t, uint64(10), row.AddState(uintState))

	// The following StringState has 1 byte per character, so the new total
	// size should be 10 + 5 + 4 = 19.
	stringState := StringState("Hello")
	assert.Equal(t, uint64(19), row.AddState(stringState))
}

func TestCannotBeWeightedMoreThanOneValue(t *testing.T) {
	state16 := Uint16State(16)
	row := NewNormalizingRow()

	row.AddState(state16)
	assert.Equal(t, float64(1), row.StateWeight(state16))
	row.AddState(state16)
	assert.Equal(t, float64(1), row.StateWeight(state16))
}

func TestNormalizingRowWeighting(t *testing.T) {
	state16 := Uint16State(16)
	state32 := Uint16State(32)
	state64 := Uint16State(64)
	row := NewNormalizingRow()

	row.AddState(state16)
	assert.Equal(t, float64(1), row.StateWeight(state16))

	row.AddState(state32)
	assert.Equal(t, float64(1)/float64(2), row.StateWeight(state16))
	assert.Equal(t, float64(1)/float64(2), row.StateWeight(state32))

	row.AddState(state64)
	assert.Equal(t, float64(1)/float64(3), row.StateWeight(state16))
	assert.Equal(t, float64(1)/float64(3), row.StateWeight(state32))
	assert.Equal(t, float64(1)/float64(3), row.StateWeight(state64))

	row.AddState(state64)
	assert.Equal(t, float64(1)/float64(4), row.StateWeight(state16))
	assert.Equal(t, float64(1)/float64(4), row.StateWeight(state32))
	assert.Equal(t, float64(2)/float64(4), row.StateWeight(state64))

	row.AddState(state32)
	assert.Equal(t, float64(1)/float64(5), row.StateWeight(state16))
	assert.Equal(t, float64(2)/float64(5), row.StateWeight(state32))
	assert.Equal(t, float64(2)/float64(5), row.StateWeight(state64))
}

func TestModelSize(t *testing.T) {
	state16 := Uint16State(16)
	state32 := Uint16State(32)
	model := NewModel()

	model.AddState(state16)
	// The start state is "start" which is 5 bytes for the key. The state16
	// is 2 bytes for the key and 4 bytes for the counter and 4 bytes for the
	// current count of the row. 5 + 2 + 4 + 4 = 15.
	assert.Equal(t, uint64(15), model.Size())

	model.AddState(state32)
	// 2 more bytes for state16 as a State key for the start state plus
	// another 2 + 4 + 4 bytes for state32 as a State key to be transitioned
	// to, its counter and the new row's total count.
	assert.Equal(t, uint64(27), model.Size())

	model.AddState(state32)
	// 2 more bytes for state32 as a State key for the start state plus
	// another 2 + 4 + 4 bytes for state32 as a State key to be transitioned
	// to, its counter and the new row's total count.
	assert.Equal(t, uint64(39), model.Size())

	model.AddState(state32)
	// No new data, as we've already seen state32 -> state32.
	assert.Equal(t, uint64(39), model.Size())

	model.AddState(state16)
	// state16 as a transition state means 2 more bytes for the key and 4
	// more bytes for the counter.
	assert.Equal(t, uint64(45), model.Size())
}
