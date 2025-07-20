package controller

//Generate by claude sonnet 4

import (
	"net/http"
	"strconv"
	"template/std"

	"github.com/gin-gonic/gin"
)

// StdController 演示 std 包功能的控制器
type StdController struct{}

// NewStdController 创建新的 StdController
func NewStdController() *StdController {
	return &StdController{}
}

// FunctionDemo 演示 std::function 功能
func (sc *StdController) FunctionDemo(c *gin.Context) {
	// 创建一个计算函数
	calculate := func(a, b float64, op string) float64 {
		switch op {
		case "add":
			return a + b
		case "multiply":
			return a * b
		case "subtract":
			return a - b
		case "divide":
			if b != 0 {
				return a / b
			}
			return 0
		default:
			return 0
		}
	}

	// 包装为 Function
	calcFn := std.NewFunction(calculate)

	// 从查询参数获取值
	aStr := c.DefaultQuery("a", "10")
	bStr := c.DefaultQuery("b", "5")
	op := c.DefaultQuery("op", "add")

	a, _ := strconv.ParseFloat(aStr, 64)
	b, _ := strconv.ParseFloat(bStr, 64)

	// 调用函数
	results := calcFn.Call(a, b, op)

	c.JSON(http.StatusOK, gin.H{
		"message": "Function demo",
		"input": gin.H{
			"a":  a,
			"b":  b,
			"op": op,
		},
		"result": results[0],
		"type":   "std::function",
	})
}

// BindDemo 演示 std::bind 功能
func (sc *StdController) BindDemo(c *gin.Context) {
	// 原始的三参数函数
	power := func(base, exponent, multiplier float64) float64 {
		result := 1.0
		for i := 0; i < int(exponent); i++ {
			result *= base
		}
		return result * multiplier
	}

	// 获取参数
	baseStr := c.DefaultQuery("base", "2")
	base, _ := strconv.ParseFloat(baseStr, 64)

	// 使用 Bind 创建平方函数（指数固定为2，乘数固定为1）
	square := std.Bind(power, std.P1, 2.0, 1.0)

	// 使用 Bind 创建立方函数
	cube := std.Bind(power, std.P1, 3.0, 1.0)

	// 使用 Bind 创建双倍平方函数
	doubleSquare := std.Bind(power, std.P1, 2.0, 2.0)

	// 调用绑定的函数
	squareResult := square.Call(base)
	cubeResult := cube.Call(base)
	doubleSquareResult := doubleSquare.Call(base)

	c.JSON(http.StatusOK, gin.H{
		"message": "Bind demo",
		"input":   base,
		"results": gin.H{
			"square":       squareResult[0],
			"cube":         cubeResult[0],
			"doubleSquare": doubleSquareResult[0],
		},
		"type": "std::bind",
	})
}

// ForwardDemo 演示 std::forward 和函数组合功能
func (sc *StdController) ForwardDemo(c *gin.Context) {
	// 定义一些简单的数学函数
	double := func(x float64) float64 {
		return x * 2
	}

	addTen := func(x float64) float64 {
		return x + 10
	}

	square := func(x float64) float64 {
		return x * x
	}

	// 获取输入值
	inputStr := c.DefaultQuery("input", "5")
	input, _ := strconv.ParseFloat(inputStr, 64)

	// 使用 Pipe 创建处理管道（从左到右）
	pipeline := std.Pipe(double, addTen, square)
	pipelineResult := pipeline.Forward(input)

	// 使用 Compose 创建组合函数（从右到左）
	composed := std.Compose(square, addTen, double)
	composedResult := composed.Forward(input)

	// 单独调用每个函数展示步骤
	step1 := double(input)
	step2 := addTen(step1)
	step3 := square(step2)

	c.JSON(http.StatusOK, gin.H{
		"message": "Forward and composition demo",
		"input":   input,
		"pipeline": gin.H{
			"description": "double -> addTen -> square",
			"result":      pipelineResult[0],
		},
		"composed": gin.H{
			"description": "square(addTen(double(x)))",
			"result":      composedResult[0],
		},
		"steps": gin.H{
			"step1_double": step1,
			"step2_addTen": step2,
			"step3_square": step3,
		},
		"type": "std::forward",
	})
}

// FunctionalDemo 演示函数式编程功能
func (sc *StdController) FunctionalDemo(c *gin.Context) {
	// 获取数组参数（从查询字符串）
	numbersParam := c.DefaultQuery("numbers", "1,2,3,4,5,6,7,8,9,10")

	// 解析数字数组
	var numbers []int
	if numbersParam != "" {
		for _, numStr := range c.QueryArray("num") {
			if num, err := strconv.Atoi(numStr); err == nil {
				numbers = append(numbers, num)
			}
		}
	}

	// 如果没有从 num 参数获取到数字，使用默认值
	if len(numbers) == 0 {
		numbers = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	}

	// Map: 所有数字乘以2
	doubled := std.Map(numbers, func(x int) int {
		return x * 2
	})

	// Filter: 只保留偶数
	evens := std.Filter(numbers, func(x int) bool {
		return x%2 == 0
	})

	// Reduce: 计算所有数字的和
	sum := std.Reduce(numbers, func(acc, x int) int {
		return acc + x
	}, 0)

	// 使用柯里化创建添加特定数字的函数
	add := func(a, b int) int {
		return a + b
	}
	curriedAdd := std.Curry2(add)
	addFive := curriedAdd(5)

	// 对所有数字加5
	addedFive := std.Map(numbers, addFive)

	c.JSON(http.StatusOK, gin.H{
		"message": "Functional programming demo",
		"input":   numbers,
		"results": gin.H{
			"doubled":   doubled,
			"evens":     evens,
			"sum":       sum,
			"addedFive": addedFive,
		},
		"type": "functional_programming",
	})
}

// AdvancedDemo 演示高级功能（记忆化、Once等）
func (sc *StdController) AdvancedDemo(c *gin.Context) {
	// 模拟一个昂贵的计算
	expensiveCalculation := func(n int) int {
		// 模拟耗时操作
		result := 0
		for i := 0; i <= n; i++ {
			result += i
		}
		return result
	}

	// 创建记忆化版本
	memoizedCalc := std.Memoize(expensiveCalculation)

	// 获取输入参数
	inputStr := c.DefaultQuery("n", "100")
	n, _ := strconv.Atoi(inputStr)

	// 调用记忆化函数
	result := memoizedCalc(n)

	// 创建一次性执行的函数
	counter := 0
	increment := func() int {
		counter++
		return counter
	}
	onceIncrement := std.Once(increment)

	// 多次调用，但只会执行一次
	call1 := onceIncrement()
	call2 := onceIncrement()
	call3 := onceIncrement()

	c.JSON(http.StatusOK, gin.H{
		"message": "Advanced features demo",
		"memoization": gin.H{
			"input":  n,
			"result": result,
		},
		"once": gin.H{
			"call1": call1,
			"call2": call2,
			"call3": call3,
			"note":  "All calls return the same value (1) because function only executes once",
		},
		"type": "advanced_features",
	})
}
