# 熟知go语言的一切


本章主要讨论go的设计哲学，并且如何基于go的设计哲学去写`go code`，能够领悟go的编程思想，并结合这些思想编写出高质量的go代码。


## go设计哲学

### 1. 简单，通过组合的方式进行解偶

go是google内部产生的一种语言，受C语言影响较大，其基本语法参考了C语言。go最显著的设计哲学就是`少即是多`，同主流面向对象语言不同的是：

go通过`组合`实现快速将新类型复用其他类型已经实现的能力：

example1: struct中嵌入struct

```go
type poolLocal struct {
	private interface{}
	shared []integace{}
	Mutex
	pad [128]byte
}
```

poolLocal类型中嵌入了Mutex，被嵌入的Mutex类型的方法集合会被提升到外面的类型中。例如，poolLocal将拥有Mutex类型的Lock和Unlock方法。但是实际调用中，
方法调用会被传给poolLocal的Mutex实例。

在go的标准库中还有很多类似的用法，通过在interface的定义中嵌入interface类型来实现接口行为的聚合，例如：

example: interface中嵌入interface

```go
type ReadWriter interface {
	Reader
	Writer
}
```

上述两种示例通过称为接口的`垂直组合`，将接口进行嵌入实现新的接口行为，演示代码：`./example/no-1/no_1.go`。go中还有一种将程序各个部分组合在一起的方法，书中称之为`水平组合`。
```go
// $GOROOT/src/io/ioutil/ioutil.go
func ReadAll(r io.Reader)([]byte, error)

// $GOROOT/src/io/io.go
func Copy(dst Writer, src Reader)(written int64, err error)
```
函数`ReadAll`通过调用`io.Reader`这个接口将`io.Reader`的实现与`ReadAll`所在的包以低耦合的方式水平组合在了一起。_后面接口章节仔细给出解释_。

### 2. 原生并发，充分利用多核计算机特性



### 3. 面向工程，标准库中实现大量直接可用的工具包
