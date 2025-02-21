package array

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

func modifyElements(s []int) {
	s[0] = 100 // 修改元素，会影响到原始切片
}

func modifySlice(s []int) {
	s = append(s, 4) // 修改切片本身（长度），不会影响到原始切片
	fmt.Println(len(s))
}

func modifySliceNew(s []int) {
	s = append(s, 4) //发生扩容
	s[0] = 100
	fmt.Println("in modifySliceNew", s)
}

func SExample() {
	// 示例 1: 修改元素
	//s1 := []int{1, 2, 3}
	//fmt.Println("Before modifyElements:", s1) // [1 2 3]
	//modifyElements(s1)
	//fmt.Println("After modifyElements:", s1) // [100 2 3] (元素被修改)
	//
	//// 示例 2: 修改切片本身
	//s2 := []int{1, 2, 3}
	//fmt.Println("Before modifySlice:", s2) // [1 2 3]
	//modifySlice(s2)
	//fmt.Println("After modifySlice:", s2, len(s2)) // [1 200 3] (长度未变，但第二个元素被修改)
	//
	//// 示例3：修改切片本身(扩容)
	//s3 := []int{1, 2, 3}
	//fmt.Println("Before modifySliceNew:", s3)
	//modifySliceNew(s3)
	//fmt.Println("After modifySliceNew", s3)

	// 示例4: 共享底层数组
	s4 := []int{1, 2, 3, 4, 5}
	fmt.Println("s4", s4)
	s5 := s4[:] // s5 共享 s4 的底层数组
	fmt.Println("s5", s5)
	s5 = append(s5, 100)
	s5[0] = 1000
	fmt.Println("After modifying s5:")
	fmt.Println("s4:", s4) // [1,2,3,4,5]
	fmt.Println("s5:", s5) // [1000, 2, 3, 4, 5,100]
}

var bigSlice []int // 全局变量，用于模拟长期持有的切片

func createBigSlice() []int {
	s := make([]int, 1000000) // 创建一个大数组
	for i := range s {
		s[i] = i
	}
	return s
}

func getSmallSlice(s []int) []int {
	return s[:10] // 返回一个小切片，但仍然引用着大数组
}

// 创建新的切片
func getSmallSlice2(s []int) []int {
	small := make([]int, 10)
	copy(small, s[:10]) // 将 s 的前 10 个元素复制到 small
	return small
}

func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func TestSlice() {
	bigSlice = createBigSlice()            // 创建一个大数组
	smallSlice := getSmallSlice2(bigSlice) // 获取一个小切片

	// 模拟使用 smallSlice，但 bigSlice 不再直接使用
	_ = smallSlice

	fmt.Println(&bigSlice[0], &smallSlice[0])

	// 强制进行垃圾回收
	runtime.GC()
	printMemStats() // 此时 bigSlice 仍然存在, 大数组无法被回收

	// 释放 bigSlice 的引用
	bigSlice = nil

	// 再次强制进行垃圾回收
	runtime.GC()
	printMemStats() // 大数组被回收
	fmt.Println(&smallSlice[0])
	fmt.Println(smallSlice)
}

func TestSlice2() {
	bigSlice2 := createBigSlice()
	smallSlice2 := getSmallSlice(bigSlice2) // getSmallSlice 返回了一个子切片
	_ = smallSlice2                         // 假装我们用了 smallSlice2
	printMemStats()

	fmt.Println("set bigSlice2 is nil")

	bigSlice2 = nil // 认为释放bigSlice2就没问题了，实际上不行

	runtime.GC()
	printMemStats() // 大数组仍然无法被回收
}

func TestSlice3() {
	// 关闭GC，以便我们能更清晰地看到内存变化
	debug.SetGCPercent(-1)

	bigSlice3 := createBigSlice()
	smallSlice := getSmallSlice(bigSlice3)
	// 打印smallSlice的地址和底层数组的起始地址.
	// 输出：smallSlice addr: 0x14000180000, bigSlice addr: 0x14000180000，由于small 和 big的指向同一底层数组，所以其打印地址相同
	fmt.Printf("smallSlice addr: %p, bigSlice addr: %p\n", &smallSlice[0], &bigSlice3[0])

	// 强制进行一次垃圾回收
	runtime.GC()
	printMemStats()

	bigSlice3 = nil

	// 等待一段时间，确保bigSlice有机会被回收（如果有的话）
	time.Sleep(time.Second)

	runtime.GC()
	printMemStats()

	// 再次打印地址, 看看smallSlice的地址是否还在原来的范围
	fmt.Printf("smallSlice addr: %p\n", &smallSlice[0])
}

func TestP() {
	s := make([]int, 0, 5)
	//s = append(s, 1, 2, 3)
	//s1 := s[1:3]
	//s2 := s[2:4]
	//s1[0] = 100
	fmt.Println(s)
	//fmt.Println(s1)
	//fmt.Println(s2)
}

func TestPP() {
	s := []int{1, 2, 3, 4, 5}
	s1 := s[1:3]
	s1 = append(s1, 10, 20, 30, 40, 50)
	fmt.Println(s)
}
