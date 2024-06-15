package blinker

import (
	"machine"
)

type LockedPWM struct {
	pwm machine.TCC
}

func (p *LockedPWM) Get() machine.TCC {
	return p.pwm
}

func (p *LockedPWM) Top() (top uint32) {
	top = p.pwm.Top()
	return
}

func (p *LockedPWM) Channel(pin machine.Pin) (ch uint8, err error) {
	ch, err = p.pwm.Channel(pin)
	return
}

func (p *LockedPWM) Set(channel uint8, value uint32) {
	p.pwm.Set(channel, value)
}

func NewLockedPWM() (*LockedPWM, error) {
	var pwm machine.TCC
	err := pwm.Configure(machine.PWMConfig{})
	return &LockedPWM{
		pwm: pwm,
	}, err
}
