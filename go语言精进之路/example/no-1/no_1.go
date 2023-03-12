package no_1

// example1：
// LGComputer中嵌入了LGIndicator，LGComputer这种实体是具有显示器的能力

// LG计算机
type LGComputer struct {
	LGIndicator
}

// LG显示器
type LGIndicator struct {
}

// example2: 还是使用计算机和显示器进行说明，不过这里的计算机和显示器都是一组抽象的概念
// Computer的方法集合是包含了显示器的方法集的
type Computer interface {
	Indicator
}

type Indicator interface {
	Display()
	Open()
	Close()
}

// 上述称为垂直组合
