package types

import "git.o0.tel/sidc/unoblink/devices"

type (
	TrinityChans [3]<-chan int
	TrinityLEDs  [3]*devices.LightEmitter
	StateStage   int8
)

func (s StateStage) String() string {
	return [...]string{"Low pins", "High pins"}[s]
}
