package array

import "fmt"

func Array() {
	var months = [13]string{
		1:  "jan",
		2:  "feb",
		3:  "mar",
		4:  "apr",
		5:  "may",
		6:  "june",
		7:  "july",
		8:  "aug",
		9:  "sep",
		10: "oct",
		11: "nov",
		12: "dec",
	}

	q2 := months[4:7]
	q3 := months[6:9]

	fmt.Println(len(q2), cap(q2))
	fmt.Println(len(q3), cap(q3))
}

func Array2() {
	var s []int
	s = nil
	s = []int(nil)
	s = []int{}

	if s == nil {
		fmt.Println("s is nil")
	}
}

/*
**1. `var s []int`**

   *   **声明:** 声明了一个名为 `s` 的切片变量，其元素类型为 `int`。
   *   **值:** `s` 的值为 `nil`。
   *   **底层数组:** 没有分配底层数组。
   *   **长度和容量:** `len(s)` 和 `cap(s)` 都为 0。
   *   **特性:**
       *   这是声明切片变量的 *零值* 方式。
       *   `s` 是一个 *nil 切片*。
       *   对 `nil` 切片进行 `append` 操作是安全的，Go 会自动分配底层数组。

**2. `s = nil`**

   *   **赋值:** 将 `nil` 赋值给切片变量 `s`。
   *   **前提:** `s` 必须已经声明过（例如，通过 `var s []int`）。
   *   **值:** `s` 的值为 `nil`。
   *   **底层数组:**  如果 `s` 之前有底层数组，这个赋值操作会解除 `s` 与其底层数组的关联（原底层数组如果没有被其他变量引用，可能会被垃圾回收）。如果没有，则无影响。
   *   **长度和容量:** `len(s)` 和 `cap(s)` 都为 0。
   * **特性**:
     * 使 `s` 成为一个 *nil 切片*.

**3. `s = []int(nil)`**

   *   **赋值:** 使用类型转换将 `nil` 转换为 `[]int` 类型，然后赋值给 `s`。
   *   **前提:**  `s` 必须已经声明过。
   *   **值:** `s` 的值为 `nil`。
   *   **底层数组:** 没有分配底层数组。
   *   **长度和容量:** `len(s)` 和 `cap(s)` 都为 0。
   * **特性**
      *   等价于 `s = nil`，将 `s` 变成一个 *nil 切片*。
      *  这种写法比较啰嗦，通常不使用。

**4. `s = []int{}`**

   *   **赋值:**  使用复合字面量创建一个空的 `int` 切片，并将其赋值给 `s`。
   *   **前提:** `s` 必须已经声明过。
   *   **值:** `s` 的值是一个 *非 nil* 的空切片。
   *   **底层数组:**  分配了一个空的底层数组（长度和容量都为 0）。
   *   **长度和容量:** `len(s)` 和 `cap(s)` 都为 0。
   * **特性:**
      * `s` 是一个 *非 nil 的空切片*。
      * 这是创建空切片的常用方法。

**关键区别：**

*   **`nil` 切片 vs. 空切片:**
    *   **`nil` 切片:** 没有底层数组。  它的值为 `nil`。
    *   **空切片:** 有一个底层数组，但数组长度为 0。它的值不是 `nil`。

*   **比较:**
    *   `nil` 切片与 `nil` 比较结果为 `true`。
    *   空切片与 `nil` 比较结果为 `false`。

**代码示例：**

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var s1 []int
	s2 := []int(nil)
	s3 := []int{}

	fmt.Println("s1 == nil:", s1 == nil)       // true
	fmt.Println("s2 == nil:", s2 == nil)       // true
	fmt.Println("s3 == nil:", s3 == nil)       // false

	fmt.Println("s1, len:", len(s1), "cap:", cap(s1)) // 0 0
	fmt.Println("s2, len:", len(s2), "cap:", cap(s2)) // 0 0
	fmt.Println("s3, len:", len(s3), "cap:", cap(s3)) // 0 0

    //通过反射进一步确认
    fmt.Println("s1", reflect.ValueOf(s1).IsNil()) //true
    fmt.Println("s2", reflect.ValueOf(s2).IsNil()) //true
    fmt.Println("s3", reflect.ValueOf(s3).IsNil()) //false

	// 对 nil 切片和空切片进行 append 操作都是安全的
	s1 = append(s1, 1)
	s3 = append(s3, 1)

	fmt.Println("s1 after append:", s1) // [1]
	fmt.Println("s3 after append:", s3) // [1]
}
```

**何时使用哪种方式：**

*   **`var s []int` 或 `s = nil`:**  当你声明一个切片变量，但暂时不需要给它分配底层数组时。  稍后可以通过 `append` 或其他方式给它赋值。
*   **`s = []int{}`:** 当你需要创建一个空的切片，并且希望它是一个非 `nil` 值时（例如，你可能需要将它传递给一个期望非 `nil` 切片的函数）。
*   **`make([]int, 0)`:** 创建一个长度为 0，容量为 0 的空切片，等价于 `[]int{}`
*  **`make([]int, 0, 10)`:** 创建一个长度为0,容量为10的切片。

通常，`var s []int` 和 `s = []int{}` 是最常用的两种方式。`s = nil` 和 `s = []int(nil)` 意义相同,但后者可读性较差。选择哪种方式取决于你的具体需求和代码风格。大多数情况下，使用`var s []int`声明，需要立即使用字面量初始化时,使用`s := []int{}`。
*/
