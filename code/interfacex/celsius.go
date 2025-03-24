package interfacex

import "fmt"

type Celsius float64

func (c Celsius) String() string { return fmt.Sprintf("%g C", c) }

type celsiusFlag struct {
	Celsius
}

func (f *celsiusFlag) Set(s string) error {
	return nil
}
