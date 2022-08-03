package acm1602ni

import (
	"fmt"
	"time"
)

type bus interface {
	ReadReg(byte, []byte) error
	WriteReg(byte, []byte) error
}

type ACM1602NI struct {
	dev  bus
	cmd  byte
	data byte
}

func New(dev bus) *ACM1602NI {
	return &ACM1602NI{
		dev:  dev,
		cmd:  0x00,
		data: 0x80,
	}
}

func delay() {
	time.Sleep(time.Millisecond * 5)
}

func (a *ACM1602NI) Init() error {
	init := []byte{0x01, 0x38, 0x0c, 0x06}
	for _, i := range init {
		err := a.dev.WriteReg(a.cmd, []byte{i})
		if err != nil {
			return fmt.Errorf("failed to initialize ACM1602NI: %v", err)
		}
		delay()
	}
	return nil
}

func (a *ACM1602NI) Show(line1, line2 string) error {
	var l1byte []byte
	for _, c := range line1 {
		l1byte = append(l1byte, byte(c))
	}
	var l2byte []byte
	for _, c := range line2 {
		l2byte = append(l2byte, byte(c))
	}

	for _, b := range l1byte {
		err := a.dev.WriteReg(a.data, []byte{b})
		if err != nil {
			return fmt.Errorf("failed to write line 1: %v", err)
		}
		delay()
	}

	err := a.dev.WriteReg(a.cmd, []byte{0x0c0})
	if err != nil {
		return fmt.Errorf("failed to write new line: %v", err)
	}
	delay()

	for _, b := range l2byte {
		err := a.dev.WriteReg(a.data, []byte{b})
		if err != nil {
			return fmt.Errorf("failed to write line 2: %v", err)
		}
		delay()
	}
	return nil
}
