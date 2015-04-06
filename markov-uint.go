package markov

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type uintType interface {
	numberOfBytes() uint8
	value() []byte
}

type Uint8 uint8

func (u Uint8) numberOfBytes() uint8 {
	return uint8(1)
}
func (u Uint8) value() []byte {
	b := make([]byte, 1)
	binary.Write(bytes.NewBuffer(b), binary.LittleEndian, u)
	return b
}

type Uint16 uint16

func (u Uint16) numberOfBytes() uint8 {
	return uint8(2)
}
func (u Uint16) value() []byte {
	b := make([]byte, 2)
	binary.Write(bytes.NewBuffer(b), binary.LittleEndian, u)
	return b
}

type Uint32 uint32

func (u Uint32) numberOfBytes() uint8 {
	return uint8(4)
}
func (u Uint32) value() []byte {
	b := make([]byte, 4)
	binary.Write(bytes.NewBuffer(b), binary.LittleEndian, u)
	return b
}

type Uint64 uint64

func (u Uint64) numberOfBytes() uint8 {
	return uint8(8)
}
func (u Uint64) value() []byte {
	b := make([]byte, 8)
	binary.Write(bytes.NewBuffer(b), binary.LittleEndian, u)
	return b
}

type uintModel struct {
	last  uintType
	size  uint64
	model map[uintType]map[uintType]uint8
}

func (u *uintModel) Size() uint64 {
	return u.size
}

func (u *uintModel) AddData(data interface{}) error {
	if sample, ok := data.(uintType); ok {
		if val, ok := u.model[u.last][sample]; ok {
			u.model[u.last][sample] = val + uint8(1)
		} else {
			u.model[u.last] = make(map[uintType]uint8)
			u.model[u.last][sample] = uint8(1)
			u.size = u.size + uint64(u.last.numberOfBytes()) +
				uint64(sample.numberOfBytes()) + uint64(1)
		}
		u.last = sample
		return nil

	}
	return errors.New(fmt.Sprintf(
		"Cannot call AddData with type %T on a UintModel.",
		data))
}
