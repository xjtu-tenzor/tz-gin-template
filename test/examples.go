package test

import (
	"fmt"
	"template/std"
)

// RunAllExamples 运行所有示例
func RunAllExamples() {
	fmt.Println("=== STD Package Examples ===")

	ExampleFunction()
	fmt.Println()

	ExampleBind()
	fmt.Println()

	ExampleForward()
	fmt.Println()

	ExampleFunctionalProgramming()
	fmt.Println()

	ExampleMemoization()
	fmt.Println()

	ExampleOnce()
	fmt.Println()
} // ExampleFunction 展示如何使用 Function
func ExampleFunction() {
	fmt.Println("=== Function Example ===")

	// 创建一个简单的加法函数
	add := func(a, b int) int {
		return a + b
	}

	// 包装为 Function
	fn := std.NewFunc2(add)

	// 调用函数
	results := fn.Call(3, 5)
	fmt.Printf("Function result: %v\n", results[0]) // 输出: 8

	// 创建一个无参数函数
	sayHello := func() string {
		return "Hello, World!"
	}

	fn0 := std.NewFunc0(sayHello)
	results0 := fn0.Call()
	fmt.Printf("No-arg function result: %v\n", results0[0])
}

// ExampleBind 展示如何使用 Bind
func ExampleBind() {
	fmt.Println("=== Bind Example ===")

	// 原始函数
	multiply := func(a, b, c int) int {
		return a * b * c
	}

	// 绑定第一个和第三个参数，保留第二个作为占位符
	bound := std.Bind(multiply, 2, std.P1, 10) // 绑定 a=2, c=10, b 使用占位符 P1

	// 调用绑定的函数，只需要提供一个参数
	result := bound.Call(5)                    // 相当于 multiply(2, 5, 10)
	fmt.Printf("Bind result: %v\n", result[0]) // 输出: 100

	// 使用便利函数
	add := func(a, b int) int {
		return a + b
	}

	// 绑定第一个参数
	addFive := std.BindFirst(add, 5)
	fmt.Printf("BindFirst result: %v\n", addFive(3)) // 输出: 8

	// 参数重排序示例
	subtract := func(a, b, c float64) float64 {
		return a - b - c
	}

	// 重新排列参数顺序
	reordered := std.Bind(subtract, std.P3, std.P1, std.P2)
	result2 := reordered.Call(1.0, 2.0, 3.0) // subtract(3, 1, 2) = 0
	fmt.Printf("Parameter reordering result: %v\n", result2[0])
}

// ExampleForward 展示如何使用 Forward
func ExampleForward() {
	fmt.Println("=== Forward Example ===")

	// 创建几个简单的函数
	double := func(x int) int {
		return x * 2
	}

	addOne := func(x int) int {
		return x + 1
	}

	square := func(x int) int {
		return x * x
	}

	// 使用 Pipe 从左到右组合函数
	pipeline := std.Pipe(double, addOne, square)

	// 调用管道: square(addOne(double(3))) = square(addOne(6)) = square(7) = 49
	result := pipeline.Forward(3)
	fmt.Printf("Pipeline result: %v\n", result[0]) // 输出: 49

	// 使用 Compose 从右到左组合函数
	composed := std.Compose(square, addOne, double)

	// 调用组合: square(addOne(double(3))) = 49
	result2 := composed.Forward(3)
	fmt.Printf("Compose result: %v\n", result2[0]) // 输出: 49
}

// ExampleFunctionalProgramming 展示函数式编程特性
func ExampleFunctionalProgramming() {
	fmt.Println("=== Functional Programming Example ===")

	// Map 示例
	numbers := []int{1, 2, 3, 4, 5}
	doubled := std.Map(numbers, func(x int) int {
		return x * 2
	})
	fmt.Printf("Map result: %v\n", doubled) // 输出: [2 4 6 8 10]

	// Filter 示例
	evens := std.Filter(numbers, func(x int) bool {
		return x%2 == 0
	})
	fmt.Printf("Filter result: %v\n", evens) // 输出: [2 4]

	// Reduce 示例
	sum := std.Reduce(numbers, func(acc, x int) int {
		return acc + x
	}, 0)
	fmt.Printf("Reduce result: %v\n", sum) // 输出: 15

	// Curry 示例
	add := func(a, b int) int {
		return a + b
	}

	curriedAdd := std.Curry2(add)
	addFive := curriedAdd(5)
	result := addFive(3)
	fmt.Printf("Curry result: %v\n", result) // 输出: 8
}

// ExampleMemoization 展示记忆化
func ExampleMemoization() {
	fmt.Println("=== Memoization Example ===")

	// 模拟一个耗时的计算
	fibonacci := func(n int) int {
		if n <= 1 {
			return n
		}
		// 这里为了演示，我们假设这是一个耗时操作
		fmt.Printf("Computing fibonacci(%d)\n", n)
		return n // 简化版本，实际应该是递归计算
	}

	// 创建记忆化版本
	memoFib := std.Memoize(fibonacci)

	// 第一次调用会执行计算
	result1 := memoFib(10)
	fmt.Printf("First call result: %v\n", result1)

	// 第二次调用会使用缓存
	result2 := memoFib(10)
	fmt.Printf("Second call result: %v\n", result2)
}

// ExampleOnce 展示只执行一次的函数
func ExampleOnce() {
	fmt.Println("=== Once Example ===")

	counter := 0

	increment := func() int {
		counter++
		fmt.Printf("Incrementing counter to %d\n", counter)
		return counter
	}

	// 创建只执行一次的版本
	onceIncrement := std.Once(increment)

	// 多次调用，但只会执行一次
	fmt.Printf("Result 1: %v\n", onceIncrement()) // 输出: 1
	fmt.Printf("Result 2: %v\n", onceIncrement()) // 输出: 1 (使用缓存)
	fmt.Printf("Result 3: %v\n", onceIncrement()) // 输出: 1 (使用缓存)
}

// ExampleComplexScenarios 展示复杂使用场景
func ExampleComplexScenarios() {
	fmt.Println("=== Complex Scenarios Example ===")

	// 场景1: 创建HTTP处理器工厂
	httpHandler := func(method, path, body string, statusCode int) string {
		return fmt.Sprintf("[%s] %s -> %s (status: %d)", method, path, body, statusCode)
	}

	// 创建专门的处理器
	getHandler := std.Bind(httpHandler, "GET", std.P1, "", 200)
	postHandler := std.Bind(httpHandler, "POST", std.P1, std.P2, 201)

	fmt.Printf("GET Handler: %s\n", getHandler.Call("/users")[0])
	fmt.Printf("POST Handler: %s\n", postHandler.Call("/users", "user data")[0])

	// 场景2: 数据处理管道
	processData := func(data []int) []int {
		// 过滤偶数 -> 乘以3 -> 过滤大于10的数
		filtered := std.Filter(data, func(x int) bool { return x%2 == 0 })
		tripled := std.Map(filtered, func(x int) int { return x * 3 })
		result := std.Filter(tripled, func(x int) bool { return x > 10 })
		return result
	}

	data := []int{1, 2, 3, 4, 5, 6, 7, 8}
	processed := processData(data)
	fmt.Printf("Data processing pipeline: %v -> %v\n", data, processed)

	// 场景3: 配置函数
	configuredLogger := func(level, message string, timestamp bool) string {
		prefix := ""
		if timestamp {
			prefix = "[2023-07-21 12:00:00] "
		}
		return fmt.Sprintf("%s[%s] %s", prefix, level, message)
	}

	infoLogger := std.Bind(configuredLogger, "INFO", std.P1, true)
	errorLogger := std.Bind(configuredLogger, "ERROR", std.P1, true)

	fmt.Printf("Info: %s\n", infoLogger.Call("Application started")[0])
	fmt.Printf("Error: %s\n", errorLogger.Call("Connection failed")[0])
}
