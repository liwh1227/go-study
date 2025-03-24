package interfacex

import "fmt"

type Speaker interface {
	Speak()
}

type Cat struct {
	Name string
}

func (c *Cat) Speak() {
	fmt.Println(c.Name, "says: Meow!")
}

func TestSpeak() {
	c := Cat{Name: "tom"}
	c.Speak()
}

type Data struct {
	Value int
}

// 值接收者方法
func (d *Data) PrintValue() {
	fmt.Println("Value:", d.Value)
}

// 指针接收者方法
func (d Data) Increment() {
	d.Value++
}

type Incrementor interface {
	Increment()
}

func TestData() {
	data := Data{Value: 10}
	dataPtr := &data

	dataPtr.PrintValue()

	data.PrintValue()

	data.Increment()

	data.PrintValue()
}

func TestData2() {
	var i Incrementor
	var dPtr = &Data{Value: 100}
	i = dPtr
	i.Increment()
	dPtr.PrintValue()
	dPtr.Increment()
	dPtr.PrintValue()
}

type MyType struct {
	Value int
}

// 值接收者方法
func (m MyType) ModifyValue() {
	m.Value = 100 // 修改的是 m 的副本，不是原始值
	fmt.Println("Inside ModifyValue:", m.Value)
}

func (m *MyType) ModifyValue2() {
	m.Value = 20
	fmt.Println("Inside ModifyValue:", m.Value)
}

func TestModifyValue() {
	original := MyType{Value: 10}
	ptr := &original

	fmt.Println("Before:", original.Value)

	original.ModifyValue()

	fmt.Println("After:", original.Value)

	ptr.ModifyValue()

	fmt.Println("After ptr modify:", original.Value)
}

func TestModifyValue2() {
	original := MyType{Value: 1000}

	original.ModifyValue2()

	fmt.Println("original value:", original.Value)
}
