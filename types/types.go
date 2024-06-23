package types

type StateStage int8

func (s StateStage) String() string {
	return [...]string{"Low pins", "High pins"}[s]
}

type PinVolt struct {
	Pin     uint8
	Voltage uint32
}
