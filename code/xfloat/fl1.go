package xfloat

import "fmt"

/*
要理解为什么 `int64` 值 3000000000 在转换为 `int32` 后变成 -1294967296，我们需要了解二进制补码表示法以及截断过程。

**1. 二进制补码 (Two's Complement):**

   *   现代计算机使用二进制补码来表示有符号整数。
   *   **正数:** 补码与原码（直接的二进制表示）相同。
   *   **负数:** 补码是其绝对值的原码按位取反（0 变 1，1 变 0），然后加 1。
   *   **最高位 (Most Significant Bit, MSB):**  在补码中，最高位是符号位。0 表示正数，1 表示负数。

**2. `int64` 值 3000000000 的二进制表示:**

   *   3000000000 的二进制表示（64 位）：

       ```
       00000000 00000000 00000000 00000000 10110010 01100100 10010100 00000000
       ```
       (为了方便阅读，每 8 位用空格隔开)

**3. 截断为 `int32`:**

   *   当将 `int64` 转换为 `int32` 时，高 32 位被丢弃，只保留低 32 位：

       ```
       10110010 01100100 10010100 00000000  (这是截断后的 32 位)
       ```

**4. `int32` 值的解释:**

   *   现在我们有一个 32 位的二进制数，需要将其解释为 `int32`（有符号 32 位整数，使用补码）。
   *   由于最高位是 1，这是一个负数。
   *   要得到负数的绝对值，我们需要对这个补码再求一次补码（按位取反，然后加 1）：
      1.  按位取反: 01001101 10011011 01101011 11111111
      2.  加1:      01001101 10011011 01101100 00000000

   *   这个二进制数 (01001101 10011011 01101100 00000000) 对应的十进制值是 1294967296。
   *   因为原始的 32 位补码表示的是负数，所以最终结果是 -1294967296。

**总结步骤：**

1.  `int64` 值 (3000000000) 用 64 位二进制表示。
2.  转换为 `int32` 时，高 32 位被截断，只保留低 32 位。
3.  截断后的 32 位二进制数，由于最高位是 1，根据补码规则，它表示一个负数。
4.  通过对该补码再次求补码（取反加一），得到其绝对值 1294967296。
5.  因此，最终的 `int32` 值为 -1294967296。

**关键点：**

*   截断会导致信息丢失。
*   补码表示法决定了负数的表示方式。
*   最高位（符号位）决定了整数的正负。

这就是为什么 3000000000 (int64) 转换为 int32 后变成 -1294967296 的原因。它演示了在进行整数类型转换时可能发生的溢出和截断，以及二进制补码表示法如何影响结果。
*/

func Float() {
	var big int64 = 3000000000 // 超出 int32 范围
	var small int32 = int32(big)

	fmt.Println(small)

	var big2 int64 = 2147483647 //int32最大值
	var small2 int32 = int32(big2)
	fmt.Println("big2:", big2)
	fmt.Println("small2:", small2) //2147483647

	var big3 int64 = 2147483648 //int32最大值 + 1
	var small3 int32 = int32(big3)
	fmt.Println("big3:", big3)
	fmt.Println("small3:", small3) //-2147483648 (截断)
}
