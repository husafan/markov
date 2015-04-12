package markov

/**
 * The markov-states.go file provides a single location to define provided
 * implementations of the markov.State interface.
 */

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

type Uint16State uint16

func (s Uint16State) Value() string {
	return strconv.FormatUint(uint64(s), 10)
}

func (s Uint16State) Size() uint64 {
	return uint64(2)
}

func (s Uint16State) Bytes() []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, uint16(s))
	return bytes
}

type StringState string

func (s StringState) Value() string {
	return string(s)
}

func (s StringState) Size() uint64 {
	return uint64(len(s))
}

func (s StringState) Bytes() []byte {
	return bytes.NewBufferString(string(s)).Bytes()
}
