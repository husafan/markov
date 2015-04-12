package markov

// A special start state.
var start State = StringState("start")

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
	last State
	data map[State]*NormalizingRow
}

func (m *Model) AddState(state State) {
	_, ok := m.data[m.last]
	if !ok {
		m.data[m.last] = NewNormalizingRow()
	}
	row := m.data[m.last]
	row.AddState(state)
	m.last = state
}

func (m *Model) Size() uint64 {
	var size uint64 = 0
	for key, value := range m.data {
		size += key.Size() + value.Size()
	}
	return size
}

func NewModel() *Model {
	return &Model{start, make(map[State]*NormalizingRow)}
}

/**
 * In a given model, the transitons between a given state and its successive
 * states is represented as set of probabilities between 0 and 1. This object
 * encapsulates the mechanics of adding transition states to a given starting
 * state and retreiving the weights associated with any given state.
 */
type NormalizingRow struct {
	// Holds the number of states that have been added to this row.
	currentCount uint32
	// Holds the counts for given state transitions.
	rowData map[State]uint32
	// Holds the current size of this row.
	size uint64
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
	value := n.rowData[state]
	// If value is 0, this state has not been seen before. Add the size of
	// the key plus the size of its counter to the current size.
	if value == 0 {
		n.size += state.Size() + 4
	}
	value = value + 1
	n.rowData[state] = value
	n.currentCount++
	return n.size
}

/**
 * Returns the weight associated with a given state, and true if the state
 * existed in the underlying data. If the state has not been seen before, this
 * method returns 0.
 */
func (n *NormalizingRow) StateWeight(state State) float64 {
	value := n.rowData[state]
	if value == 0 {
		return float64(value)
	}
	return float64(value) / float64(n.currentCount)
}

/**
 * Returns a new NormalizingRow initialized to have no values, a current
 * weight equal to 1 and a size accounting only for the weight value.
 */
func NewNormalizingRow() *NormalizingRow {
	return &NormalizingRow{0, make(map[State]uint32), uint64(4)}
}
