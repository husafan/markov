package markov

import (
	"regexp"
	"testing"

	"github.com/husafan/markov"
	"github.com/stretchr/testify/assert"
)

func TestUintModelCanOnlyAcceptUintTypes(t *testing.T) {
	model := markov.NewUintModel()
	assert.NotNil(t, model)

	var err error
	var re *regexp.Regexp
	err = model.AddData("Hello")
	assert.NotNil(t, err)
	re = regexp.MustCompile("type string")
	assert.NotEqual(t, "", re.FindString(err.Error()))

	err = model.AddData(1)
	assert.NotNil(t, err)
	re = regexp.MustCompile("type int")
	assert.NotEqual(t, "", re.FindString(err.Error()))
}

func TestUintModelSizeCalculations(t *testing.T) {
	model := markov.NewUintModel()

	var err error
	err = model.AddData(markov.Uint16(32))
	assert.Nil(t, err)
	// The model should have a single entry from 0 to 32 with a count of 1.
	// This means 1 byte for the 0, 2 bytes for the uint16 and 1 byte for
	// the count of 1.
	assert.Equal(t, uint64(4), model.Size())

	err = model.AddData(markov.Uint16(32))
	assert.Nil(t, err)
	// Two new entries of size uint16 and a new count of size uint8 adds
	// another 5 bytes to the size.
	assert.Equal(t, uint64(9), model.Size())

	err = model.AddData(markov.Uint16(32))
	assert.Nil(t, err)
	// This entry simply increments the count, so no added size.
	assert.Equal(t, uint64(9), model.Size())
}
