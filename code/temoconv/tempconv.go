// cf 将数值参数转换为摄氏温度和华氏温度

package temoconv

import "fmt"

type Celsius float64
type Fahrenheit float64

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC     Celsius = 0
	BoilingC      Celsius = 100
)

func (c Celsius) String() string { return fmt.Sprintf("%g C", c) }

func (f Fahrenheit) String() string { return fmt.Sprintf("%g C", f) }
