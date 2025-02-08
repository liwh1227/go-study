package gc

import (
	"fmt"
	"runtime"
	"time"
)

var globalCache map[string]*Data

type Data struct {
	Value     string
	LargeData [1024 * 1024]byte // 1MB 数组
}

func processRequest(request string) {
	data := &Data{Value: request}
	globalCache[request] = data
	fmt.Printf("Processing: %s,  Cache Size: %d\n", data.Value, len(globalCache)) //打印缓存大小.
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func Ff() {
	globalCache = make(map[string]*Data)

	for i := 0; i < 10000; i++ { // 增加循环次数
		processRequest(fmt.Sprintf("Request %d", i))
		if i%1000 == 0 { // 每 1000 次迭代打印一次
			printMemUsage()
		}
		time.Sleep(time.Microsecond) // 稍微减慢速度
	}

	// 模拟内存泄漏的情况
	fmt.Println("---------- After loop (with memory leak) ----------")
	runtime.GC() // 循环结束后手动GC一次，观察泄漏后的内存占用
	printMemUsage()

	// 修复内存泄漏
	fmt.Println("---------- After fixing memory leak ----------")
	globalCache = nil
	runtime.GC()
	printMemUsage()
}
