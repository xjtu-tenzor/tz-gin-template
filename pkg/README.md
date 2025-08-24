# 生成自我的AI

# pkg Package - C++ STL-like functionality for Go

这个包提供了类似 C++ STL 的功能，包括 `std::function`、`std::bind`等特性的 Go 语言实现。

## 功能特性

### 1. Function (std::function)
类似 C++ 的 `std::function`，用于存储和调用任意函数。

```go
// 创建函数包装器
add := func(a, b int) int {
    return a + b
}

fn := pkg.NewFunction(add)
results := fn.Call(3, 5) // [8]
```

### 2. Bind (std::bind)
类似 C++ 的 `std::bind`，用于绑定函数参数。

#### 基本用法
```go
multiply := func(a, b, c int) int {
    return a * b * c
}

// 绑定第一个和第三个参数，保留第二个作为占位符
bound := pkg.Bind(multiply, 2, pkg.P1, 10)
result := bound.Call(5) // multiply(2, 5, 10) = 100
```

#### 占位符详解
占位符 `pkg.P1`, `pkg.P2`, `pkg.P3` 等用于指定调用时动态提供的参数位置：

```go
// 原始函数: calculate(a, b, c, d, e) = (a + b) * c - d + e
calculate := func(a, b, c, d, e float64) float64 {
    return (a + b) * c - d + e
}

// 示例1: 绑定中间参数
bound1 := pkg.Bind(calculate, pkg.P1, pkg.P2, 3.0, 1.0, pkg.P3)
// 调用时: bound1.Call(2, 4, 5) => (2+4)*3-1+5 = 22

// 示例2: 参数重排序
bound2 := pkg.Bind(calculate, pkg.P3, pkg.P1, pkg.P2, 0.0, 10.0)
// 调用时: bound2.Call(2, 3, 1) => (1+2)*3-0+10 = 19

// 示例3: 创建专用函数
multiply3 := func(a, b, c float64) float64 { return a * b * c }
square := pkg.Bind(multiply3, pkg.P1, pkg.P1, 1.0)  // x^2
cube := pkg.Bind(multiply3, pkg.P1, pkg.P1, pkg.P1)  // x^3
```


### 3. 函数式编程特性

```go
numbers := []int{1, 2, 3, 4, 5}

// Map
doubled := pkg.Map(numbers, func(x int) int { return x * 2 })

// Filter
evens := pkg.Filter(numbers, func(x int) bool { return x%2 == 0 })

// Reduce
sum := pkg.Reduce(numbers, func(acc, x int) int { return acc + x }, 0)

// Curry
add := func(a, b int) int { return a + b }
curriedAdd := pkg.Curry2(add)
addFive := curriedAdd(5)
result := addFive(3) // 8
```

### 4. 高级功能

```go
// 记忆化
fibonacci := func(n int) int { /* 计算斐波那契数列 */ }
memoFib := pkg.Memoize(fibonacci)

// 只执行一次
counter := 0
increment := func() int { counter++; return counter }
onceIncrement := pkg.Once(increment)
```

## 占位符

在使用 `Bind` 时，可以使用以下占位符：

- `pkg.P1` - 第一个参数的占位符
- `pkg.P2` - 第二个参数的占位符
- `pkg.P3` - 第三个参数的占位符
- ... 依此类推到 `pkg.P10`

## 使用示例

```go
package main

import (
    "fmt"
    "template/pkg"
)

func main() {
    // 1. 使用 Function
    multiply := func(a, b int) int {
        return a * b
    }
    
    fn := pkg.NewFunction(multiply)
    result := fn.Call(6, 7)
    fmt.Println("Function result:", result[0]) // 42

    // 2. 使用 Bind
    power := func(base, exp, multiplier int) int {
        result := 1
        for i := 0; i < exp; i++ {
            result *= base
        }
        return result * multiplier
    }
    
    square := pkg.Bind(power, pkg._1, 2, 1)
    fmt.Println("Square of 5:", square.Call(5)[0]) // 25

}
```

## 注意事项

1. 由于 Go 的类型系统限制，某些高级特性可能不如 C++ STL 那样灵活
2. 性能可能会受到反射的影响，在性能敏感的代码中请谨慎使用
3. 这个包主要用于演示和学习目的，生产环境使用请进行充分测试

## 扩展

您可以根据需要扩展这个包，添加更多的 STL 风格的功能，如：

- 容器类型（vector, list 等）
- 算法函数（sort, find 等）
- 迭代器模式
- 智能指针模拟
- 等等...
