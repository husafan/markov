package markov

import (
	"errors"
	"fmt"
)

// A special start state.
var Start State = StringState("start")

/**
 * The State interface interface must be implemented by all types of states
 * added to a Markov model. Note that implementations of this interface will
 * serve as keys into the model's underlying data structure.
 */
type State interface {
	Value() string
	Size() uint64
	Bytes() []byte
}

/**
 * Markov models are traditionally build by successively adding a given type
 * of state and recording how often one type appears after another.
 */
type Model struct {
	current State
	data    map[State]*NormalizingRow
}

/**
 * Adds a State to the current model. The model will record a transition from
 * its current State to the State being added.
 */
func (m *Model) AddState(state State) {
	_, ok := m.data[m.current]
	if !ok {
		m.data[m.current] = NewNormalizingRow()
	}
	row := m.data[m.current]
	row.AddState(state)
	m.current = state
}

/**
 * Returns the size, in bytes, of the current model.
 */
func (m *Model) Size() uint64 {
	var size uint64 = 0
	for key, value := range m.data {
		size += key.Size() + value.Size()
	}
	return size
}

/**
 * Sets the current state of the model. Returns an error if the state has not
 * been seen before and nil otherwise.
 */
func (m *Model) SetCurrentState(state State) error {
	_, ok := m.data[state]
	if !ok {
		return errors.New(fmt.Sprintf(
			"State %v cannot be used as a starting state.", state))
	}
	m.current = state
	return nil
}

func NewModel() *Model {
	return &Model{Start, make(map[State]*NormalizingRow)}
}

/**
 * In a given model, the transitons between a given state and its successive
 * states is represented as set of probabilities between 0 and 1. This object
 * encapsulates the mechanics of adding transition states to a given starting
 * state and retreiving the weights associated with any given state.
 */
type NormalizingRow struct {
	count        uint32
	data         map[State]uint32
	size         uint64
	dirty        bool
	weightValues []float64
	weightStates []State
}

/**
 * Prepares this row for walking. It must be called before the row can be
 * walked. This method populates weightValues and weightStates so that they can
 * be used for generating next states based on a float value between 0 and 1.
 */
func (n *NormalizingRow) prepare() {
	count := len(n.data)
	n.weightValues = make([]float64, count)
	n.weightStates = make([]State, count)
	cutoff := float64(0)
	for key := range n.data {
		n.weightStates = append(n.weightStates, key)
		cutoff += n.StateWeight(key)
		n.weightValues = append(n.weightValues, cutoff)
	}
	n.dirty = false
}

/**
 * Returns the current size of the row in bytes.
 */
func (n *NormalizingRow) Size() uint64 {
	return n.size
}

/**
 * Adds the occurence of a given state to the row. This method performs any
 * renormalization that might need to occur as a result of the new value
 * being added. It returns the new current size of the row in bytes.
 */
func (n *NormalizingRow) AddState(state State) uint64 {
	value := n.data[state]
	// If value is 0, this state has not been seen before. Add the size of
	// the key plus the size of its counter to the current size.
	if value == 0 {
		n.size += state.Size() + 4
	}
	value = value + 1
	n.data[state] = value
	n.count++
	n.dirty = true
	return n.size
}

/**
 * Returns the weight associated with a given state, and true if the state
 * existed in the underlying data. If the state has not been seen before, this
 * method returns 0.
 */
func (n *NormalizingRow) StateWeight(state State) float64 {
	value := n.data[state]
	if value == 0 {
		return float64(value)
	}
	return float64(value) / float64(n.count)
}

/**
 * Walks the current model given a probability between 0 and 1. Based on the
 * sampled probability of the transition States stored in this row, this method
 * will return a State that correlates with the passed in float. If there is an
 * error, the State will be nil and an error will be returned.
 */
func (n *NormalizingRow) Walk(value float64) (State, error) {
	if value < 0 || value > 1 {
		return nil, errors.New(
			"The passed in probability must be between 0 and 1")
	}
	if len(n.data) == 0 {
		return nil, errors.New("Cannot walk a row with no data.")
	}
	if n.dirty {
		n.prepare()
	}
	for index, limit := range n.weightValues {
		if value <= limit {
			return n.weightStates[index], nil
		}
	}
	return nil, errors.New("Unexpected Walk error. No State found.")
}

/**
 * Returns a new NormalizingRow initialized to have no values, a current
 * weight equal to 1 and a size accounting only for the weight value.
 */
func NewNormalizingRow() *NormalizingRow {
	return &NormalizingRow{
		count: 0,
		data:  make(map[State]uint32),
		size:  uint64(4),
		dirty: true,
	}
}
