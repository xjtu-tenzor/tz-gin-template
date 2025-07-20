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

	ExampleCombinedFeatures()
	fmt.Println()
}

// ExampleFunction 展示如何使用 Function
func ExampleFunction() {
	fmt.Println("=== Function Example ===")

	// 创建一个简单的加法函数
	add := func(a, b int) int {
		return a + b
	}

	// 包装为 Function
	fn := std.NewFunction(add)

	// 调用函数
	results := fn.Call(3, 5)
	fmt.Printf("Function result: %v\n", results[0]) // 输出: 8

	// 创建一个无参数函数
	sayHello := func() string {
		return "Hello, World!"
	}

	fn0 := std.NewFunction(sayHello) // 使用 NewFunction 包装无参数函数
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

// ExampleCombinedFeatures 展示 Function、Bind、Forward 结合使用
func ExampleCombinedFeatures() {
	fmt.Println("=== Combined Features Example ===")

	// 场景1: 数学计算管道 - Function + Forward
	fmt.Println("场景1: 数学计算管道")

	// 定义基础数学函数
	power := func(base, exponent float64) float64 {
		result := 1.0
		for i := 0; i < int(exponent); i++ {
			result *= base
		}
		return result
	}

	add := func(a, b float64) float64 { return a + b }
	multiply := func(a, b float64) float64 { return a * b }

	// 使用 Bind 创建专门的函数
	square := std.Bind(power, std.P1, 2.0)    // 平方函数: x^2
	addTen := std.Bind(add, std.P1, 10.0)     // 加10: x + 10
	double := std.Bind(multiply, std.P1, 2.0) // 乘2: x * 2

	// 使用 Forward 创建复合运算：double(addTen(square(x)))
	squareFn := std.NewForwardingFunction(square.Call)
	pipeline := squareFn.Chain(addTen.Call).Chain(double.Call)

	result1 := pipeline.Forward(3.0) // square(3) = 9, addTen(9) = 19, double(19) = 38
	fmt.Printf("  数学管道 square(3) -> +10 -> *2 = %v\n", result1[0])

	// 场景2: 文本处理工厂 - Bind + Function + Pipe
	fmt.Println("\n场景2: 文本处理工厂")

	formatText := func(prefix, text, suffix string, uppercase bool) string {
		result := prefix + text + suffix
		if uppercase {
			result = fmt.Sprintf("%s", result) // 简化版大写转换
		}
		return result
	}

	// 创建不同的文本格式化器
	htmlBold := std.Bind(formatText, "<b>", std.P1, "</b>", false)
	logEntry := std.Bind(formatText, "[LOG] ", std.P1, " [END]", true)

	// 演示 SQL 引用格式化器的用法
	sqlQuote := std.Bind(formatText, "'", std.P1, "'", false)
	sqlResult := sqlQuote.Call("user_name")
	fmt.Printf("  SQL格式化: %v\n", sqlResult[0])

	// 使用 Function 包装后组合
	boldFn := std.NewFunction(htmlBold.Call)
	logFn := std.NewFunction(logEntry.Call)

	// 创建复合文本处理器：先加粗，再包装为日志格式
	textProcessor := std.Pipe(
		func(text string) []interface{} { return boldFn.Call(text) },
		func(result []interface{}) []interface{} { return logFn.Call(result[0]) },
	)

	processedText := textProcessor.Forward("Hello World")
	fmt.Printf("  文本处理: %v\n", processedText[0])

	// 场景3: HTTP 请求构建器 - 三者完全结合
	fmt.Println("\n场景3: HTTP 请求构建器")

	buildRequest := func(method, host, path, query, body string, port int) string {
		url := fmt.Sprintf("%s:%d%s", host, port, path)
		if query != "" {
			url += "?" + query
		}
		return fmt.Sprintf("%s %s\nBody: %s", method, url, body)
	}

	// 创建不同类型的请求构建器
	getBuilder := std.Bind(buildRequest, "GET", std.P1, std.P2, std.P3, "", 80)
	postBuilder := std.Bind(buildRequest, "POST", std.P1, std.P2, "", std.P3, 443)

	// 使用 Function 包装
	getFn := std.NewFunction(getBuilder.Call)

	// 演示 POST 请求构建
	postFn := std.NewFunction(postBuilder.Call)
	postExample := postFn.Call("api.example.com", "/api/data", `{"key": "value"}`)
	fmt.Printf("  POST请求示例: %v\n", postExample[0])

	// 创建请求验证和日志记录的管道
	logRequest := func(request string) string {
		return fmt.Sprintf("[REQUEST] %s [/REQUEST]", request)
	}

	validateRequest := func(request string) string {
		if len(request) < 10 {
			return "[INVALID] " + request
		}
		return "[VALID] " + request
	}

	// 组合：构建请求 -> 验证 -> 记录日志
	requestPipeline := std.Pipe(
		func(args ...interface{}) interface{} {
			return getFn.Call(args...)[0]
		},
		validateRequest,
		logRequest,
	)

	finalRequest := requestPipeline.Forward("example.com", "/api/users", "limit=10")
	fmt.Printf("  HTTP请求处理:\n%v\n", finalRequest[0])

	// 场景4: 数据转换链 - 展示复杂的组合
	fmt.Println("\n场景4: 数据转换链")

	// 模拟数据处理函数
	parseData := func(raw string, format string) []int {
		// 简化的数据解析
		if format == "csv" && raw == "1,2,3,4,5" {
			return []int{1, 2, 3, 4, 5}
		}
		return []int{}
	}

	transformData := func(data []int, multiplier int, offset int) []int {
		result := make([]int, len(data))
		for i, v := range data {
			result[i] = v*multiplier + offset
		}
		return result
	}

	summarizeData := func(data []int) string {
		if len(data) == 0 {
			return "No data"
		}
		sum := 0
		for _, v := range data {
			sum += v
		}
		return fmt.Sprintf("Count: %d, Sum: %d, Average: %.1f",
			len(data), sum, float64(sum)/float64(len(data)))
	}

	// 创建绑定的处理器
	csvParser := std.Bind(parseData, std.P1, "csv")
	scaleAndShift := std.Bind(transformData, std.P1, 2, 10) // 乘2加10

	// 组合完整的数据处理管道
	dataProcessor := std.Pipe(
		func(raw string) interface{} { return csvParser.Call(raw)[0] },
		func(data interface{}) interface{} { return scaleAndShift.Call(data)[0] },
		func(data interface{}) interface{} { return summarizeData(data.([]int)) },
	)

	summary := dataProcessor.Forward("1,2,3,4,5")
	fmt.Printf("  数据处理结果: %v\n", summary[0])

	// 场景5: 函数式配置模式
	fmt.Println("\n场景5: 函数式配置模式")

	createConfig := func(env, db, cache, logger string, debug bool) string {
		config := fmt.Sprintf("Environment: %s\nDatabase: %s\nCache: %s\nLogger: %s\nDebug: %v",
			env, db, cache, logger, debug)
		return config
	}

	// 创建不同环境的配置生成器
	devConfig := std.Bind(createConfig, "development", std.P1, std.P2, "console", true)
	prodConfig := std.Bind(createConfig, "production", std.P1, std.P2, "file", false)

	// 演示生产环境配置
	prodSettings := prodConfig.Call("postgresql://prod", "redis://prod")
	prodStr := prodSettings[0].(string)
	if len(prodStr) > 50 {
		fmt.Printf("  生产环境配置预览: %s...\n", prodStr[:50])
	} else {
		fmt.Printf("  生产环境配置预览: %s\n", prodStr)
	}

	// 创建配置验证和美化的管道
	beautifyConfig := func(config string) string {
		return "=== Configuration ===\n" + config + "\n====================="
	}

	configPipeline := std.NewForwardingFunction(devConfig.Call).
		Chain(beautifyConfig)

	devSettings := configPipeline.Forward("postgresql://dev", "redis://dev")
	fmt.Printf("  开发环境配置:\n%v\n", devSettings[0])

	fmt.Println("\n✅ Function、Bind、Forward 结合使用测试完成！")
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
