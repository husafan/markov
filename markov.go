package markov

// Markov models are traditionally build by successfively adding a given type
// of data and recording how often one type appears after another.
type Model interface {
	AddData(interface{}) error
	Size() uint64
}

func NewUintModel() Model {
	return &uintModel{
		last:  Uint8(0),
		size:  uint64(0),
		model: make(map[uintType]map[uintType]uint8),
	}
}
