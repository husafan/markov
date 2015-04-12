package markov

// A special start state.
var start State = StringState("start")

// A constant that represents the value of 1.
const oneValue uint32 = 1000000000

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
	size uint64
	data map[State]*NormalizingRow
}

func (m *Model) AddState(state *State) {
	_, ok := m.data[m.last]
	if !ok {
		m.data[m.last] = NewNormalizingRow()
	}
}

func (m *Model) Size() uint64 {
	return uint64(0)
}

func NewModel() *Model {
	return &Model{start, uint64(0), make(map[State]*NormalizingRow)}
}

/**
 * In a given model, the transitons between a given state and its successive
 * states is represented as set of probabilities between 0 and 1. This object
 * encapsulates the mechanics of adding transition states to a given starting
 * state and re-balancing the weights so that the probabilities remain
 * normalized.
 */
type NormalizingRow struct {
	// Holds the weight of the next value added to this row. Beause the
	// weight is between 0 and 1, this value represents the digits after
	// the decimal.
	currentWeight uint32
	// Holds the transition probabilities.
	rowData map[State]uint32
	// Holds the current size of this row.
	size uint64
}

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
	value = value + n.currentWeight
	n.rowData[state] = value
	n.size = uint64(4)
	for key, value := range n.rowData {
		if key != state {
			n.rowData[key] = value - n.currentWeight
		}
		n.size = n.size + key.Size()
	}
	n.currentWeight = n.currentWeight / 2
	return n.size
}

/**
 * Returns the weight associated with a given state, and true if the state
 * existed in the underlying data. If the state has not been seen before, this
 * method returns 0 and false.
 */
func (n *NormalizingRow) StateWeight(state State) (uint32, bool) {
	value, ok := n.rowData[state]
	return value, ok
}

/**
 * Returns a new NormalizingRow initialized to have no values, a current
 * weight equal to 1 and a size accounting only for the weight value.
 */
func NewNormalizingRow() *NormalizingRow {
	return &NormalizingRow{oneValue, make(map[State]uint32), uint64(4)}
}
