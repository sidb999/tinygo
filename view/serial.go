package view

import (
	"fmt"
	"machine"

	"git.o0.tel/sidc/unoblink/types"
)

const (
	LowP types.StateStage = iota
	HighP
)
const ADCName = "LED"

func L(n int) string {
	return fmt.Sprintf("%s%s", ADCName, string(n))
}

func GetVoltage(led machine.ADC) string {
	// vv := float32(led.Get()) * 3.3 * 1.1 * 2 / 65536 // 1.1 is just a kluge to get closer to expected battery at full charge
	// return fmt.Sprintf("%.2f (%x)", vv, v)
	return fmt.Sprintf("%.2f", float32(led.Get())*3.3/65536)
}

//	func PrintStats(stage StateStage, leds ...*devices.LightEmitter) {
//		out := "\n" + stage.String()
//		for n, led := range leds {
//			out += "\n\t" + L(n) + ": " + GetVoltage(led.ADC)
//		}
//		out += "\n"
//		fmt.Println(out)
//	}
func PrintStats(leds types.TrinityLEDs) {
	out := "\n"
	for n, led := range leds {
		out += "\n\t" + L(n) + ": " + GetVoltage(led.ADC)
	}
	out += "\n"
	fmt.Println(out)
}
