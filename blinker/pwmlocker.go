package blinker

import (
	"machine"
	"sync"
)

type LockedPWM struct {
	lock sync.Mutex
	pwm  machine.PWM
}

func (p *LockedPWM) Get() machine.PWM {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.pwm
}

func (p *LockedPWM) Top() (top uint32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	top = p.pwm.Top()
	return
}

func (p *LockedPWM) Channel(pin machine.Pin) (ch uint8, err error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	ch, err = p.pwm.Channel(pin)
	return
}

func (p *LockedPWM) Set(channel uint8, value uint32) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.pwm.Set(channel, value)
}

func NewLockedPWM() *LockedPWM {
	var pwm machine.PWM
	pwm.Configure(machine.PWMConfig{})
	return &LockedPWM{
		pwm: pwm,
	}
}
