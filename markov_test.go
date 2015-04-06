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
