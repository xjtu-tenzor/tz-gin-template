package test

import (
	"template/std"
	"testing"
)

// TestFunction 测试 std::function 功能
func TestFunction(t *testing.T) {
	t.Run("Test Function with 2 parameters", func(t *testing.T) {
		add := func(a, b int) int {
			return a + b
		}

		fn := std.NewFunc2(add)
		results := fn.Call(3, 5)

		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}

		if results[0] != 8 {
			t.Errorf("Expected 8, got %v", results[0])
		}
	})

	t.Run("Test Function with 0 parameters", func(t *testing.T) {
		getValue := func() int {
			return 42
		}

		fn := std.NewFunc0(getValue)
		results := fn.Call()

		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}

		if results[0] != 42 {
			t.Errorf("Expected 42, got %v", results[0])
		}
	})

	t.Run("Test Function with 3 parameters", func(t *testing.T) {
		calculate := func(a, b float64, op string) float64 {
			switch op {
			case "add":
				return a + b
			case "multiply":
				return a * b
			default:
				return 0
			}
		}

		fn := std.NewFunc3(calculate)
		results := fn.Call(5.0, 3.0, "add")

		if results[0] != 8.0 {
			t.Errorf("Expected 8.0, got %v", results[0])
		}

		results2 := fn.Call(4.0, 6.0, "multiply")
		if results2[0] != 24.0 {
			t.Errorf("Expected 24.0, got %v", results2[0])
		}
	})
}

// TestBind 测试 std::bind 功能
func TestBind(t *testing.T) {
	t.Run("Test basic binding", func(t *testing.T) {
		multiply := func(a, b, c int) int {
			return a * b * c
		}

		bound := std.Bind(multiply, 2, std.P1, 10)
		result := bound.Call(5)

		if result[0] != 100 {
			t.Errorf("Expected 100, got %v", result[0])
		}
	})

	t.Run("Test parameter reordering", func(t *testing.T) {
		subtract := func(a, b, c float64) float64 {
			return a - b - c
		}

		// 重新排列参数顺序: P3, P1, P2
		reordered := std.Bind(subtract, std.P3, std.P1, std.P2)
		result := reordered.Call(1.0, 2.0, 3.0) // subtract(3, 1, 2) = 3-1-2 = 0

		expected := 0.0
		actual := result[0].(float64)

		// 由于浮点数精度，我们使用接近比较
		if actual != expected {
			t.Errorf("Expected %.1f, got %.1f", expected, actual)
		}
	})

	t.Run("Test multiple placeholders", func(t *testing.T) {
		calculate := func(a, b, c, d float64) float64 {
			return (a+b)*c - d
		}

		bound := std.Bind(calculate, std.P1, 10.0, std.P2, 5.0)
		result := bound.Call(3.0, 4.0) // (3 + 10) * 4 - 5 = 47

		if result[0] != 47.0 {
			t.Errorf("Expected 47.0, got %v", result[0])
		}
	})

	t.Run("Test convenience functions", func(t *testing.T) {
		add := func(a, b int) int {
			return a + b
		}

		addFive := std.BindFirst(add, 5)
		result := addFive(3)

		if result != 8 {
			t.Errorf("Expected 8, got %v", result)
		}

		addToTen := std.BindSecond(add, 10)
		result2 := addToTen(7)

		if result2 != 17 {
			t.Errorf("Expected 17, got %v", result2)
		}
	})
}

// TestForward 测试 std::forward 功能
func TestForward(t *testing.T) {
	t.Run("Test basic forwarding", func(t *testing.T) {
		double := func(x int) int {
			return x * 2
		}

		ff := std.NewForwardingFunction(double)
		result := ff.Forward(5)

		if result[0] != 10 {
			t.Errorf("Expected 10, got %v", result[0])
		}
	})

	t.Run("Test function chaining", func(t *testing.T) {
		double := func(x int) int {
			return x * 2
		}

		addOne := func(x int) int {
			return x + 1
		}

		ff1 := std.NewForwardingFunction(double)
		chained := ff1.Chain(addOne)

		result := chained.Forward(3) // addOne(double(3)) = addOne(6) = 7

		if len(result) != 1 {
			t.Errorf("Expected 1 result, got %d", len(result))
			return
		}

		expected := 7
		actual := result[0].(int)

		if actual != expected {
			t.Errorf("Expected %d, got %d", expected, actual)
		}
	})

	t.Run("Test Pipe", func(t *testing.T) {
		double := func(x int) int {
			return x * 2
		}

		addOne := func(x int) int {
			return x + 1
		}

		square := func(x int) int {
			return x * x
		}

		pipeline := std.Pipe(double, addOne, square)
		result := pipeline.Forward(3) // square(addOne(double(3))) = square(7) = 49

		if len(result) != 1 {
			t.Errorf("Expected 1 result, got %d", len(result))
			return
		}

		expected := 49
		actual := result[0].(int)

		if actual != expected {
			t.Errorf("Expected %d, got %d", expected, actual)
		}
	})

	t.Run("Test Compose", func(t *testing.T) {
		double := func(x int) int {
			return x * 2
		}

		addOne := func(x int) int {
			return x + 1
		}

		square := func(x int) int {
			return x * x
		}

		composed := std.Compose(square, addOne, double)
		result := composed.Forward(3) // square(addOne(double(3))) = 49

		if len(result) != 1 {
			t.Errorf("Expected 1 result, got %d", len(result))
			return
		}

		expected := 49
		actual := result[0].(int)

		if actual != expected {
			t.Errorf("Expected %d, got %d", expected, actual)
		}
	})
}

// TestFunctionalProgramming 测试函数式编程特性
func TestFunctionalProgramming(t *testing.T) {
	t.Run("Test Map", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		doubled := std.Map(numbers, func(x int) int {
			return x * 2
		})

		expected := []int{2, 4, 6, 8, 10}
		for i, v := range doubled {
			if v != expected[i] {
				t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
			}
		}
	})

	t.Run("Test Filter", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6}
		evens := std.Filter(numbers, func(x int) bool {
			return x%2 == 0
		})

		expected := []int{2, 4, 6}
		if len(evens) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(evens))
		}

		for i, v := range evens {
			if v != expected[i] {
				t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
			}
		}
	})

	t.Run("Test Reduce", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		sum := std.Reduce(numbers, func(acc, x int) int {
			return acc + x
		}, 0)

		if sum != 15 {
			t.Errorf("Expected 15, got %d", sum)
		}

		product := std.Reduce(numbers, func(acc, x int) int {
			return acc * x
		}, 1)

		if product != 120 {
			t.Errorf("Expected 120, got %d", product)
		}
	})

	t.Run("Test Curry", func(t *testing.T) {
		add := func(a, b int) int {
			return a + b
		}

		curriedAdd := std.Curry2(add)
		addFive := curriedAdd(5)
		result := addFive(3)

		if result != 8 {
			t.Errorf("Expected 8, got %d", result)
		}

		multiply := func(a, b, c int) int {
			return a * b * c
		}

		curriedMul := std.Curry3(multiply)
		mulByTwo := curriedMul(2)
		mulByTwoAndThree := mulByTwo(3)
		result2 := mulByTwoAndThree(4) // 2 * 3 * 4 = 24

		if result2 != 24 {
			t.Errorf("Expected 24, got %d", result2)
		}
	})

	t.Run("Test Partial", func(t *testing.T) {
		add := func(a, b int) int {
			return a + b
		}

		addTen := std.Partial(add, 10)
		result := addTen(5)

		if result != 15 {
			t.Errorf("Expected 15, got %d", result)
		}
	})
}

// TestAdvancedFeatures 测试高级功能
func TestAdvancedFeatures(t *testing.T) {
	t.Run("Test Memoize", func(t *testing.T) {
		callCount := 0
		expensiveFunc := func(n int) int {
			callCount++
			return n * n
		}

		memoized := std.Memoize(expensiveFunc)

		// 第一次调用
		result1 := memoized(5)
		if result1 != 25 {
			t.Errorf("Expected 25, got %d", result1)
		}
		if callCount != 1 {
			t.Errorf("Expected 1 call, got %d", callCount)
		}

		// 第二次调用相同参数，应该使用缓存
		result2 := memoized(5)
		if result2 != 25 {
			t.Errorf("Expected 25, got %d", result2)
		}
		if callCount != 1 {
			t.Errorf("Expected 1 call (cached), got %d", callCount)
		}

		// 不同参数，应该重新计算
		result3 := memoized(6)
		if result3 != 36 {
			t.Errorf("Expected 36, got %d", result3)
		}
		if callCount != 2 {
			t.Errorf("Expected 2 calls, got %d", callCount)
		}
	})

	t.Run("Test Once", func(t *testing.T) {
		callCount := 0
		increment := func() int {
			callCount++
			return callCount
		}

		once := std.Once(increment)

		// 多次调用，但只会执行一次
		result1 := once()
		result2 := once()
		result3 := once()

		if result1 != 1 {
			t.Errorf("Expected 1, got %d", result1)
		}
		if result2 != 1 {
			t.Errorf("Expected 1, got %d", result2)
		}
		if result3 != 1 {
			t.Errorf("Expected 1, got %d", result3)
		}
		if callCount != 1 {
			t.Errorf("Expected function to be called only once, got %d calls", callCount)
		}
	})

	t.Run("Test Invoke", func(t *testing.T) {
		add := func(a, b int) int {
			return a + b
		}

		results := std.Invoke(add, 3, 7)
		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}
		if results[0] != 10 {
			t.Errorf("Expected 10, got %v", results[0])
		}
	})
}

// TestRef 测试引用功能
func TestRef(t *testing.T) {
	t.Run("Test Ref operations", func(t *testing.T) {
		value := 42
		ref := std.NewRef(&value)

		// 测试 Get
		if ref.Get() != 42 {
			t.Errorf("Expected 42, got %d", ref.Get())
		}

		// 测试 Set
		ref.Set(100)
		if value != 100 {
			t.Errorf("Expected original value to be 100, got %d", value)
		}
		if ref.Get() != 100 {
			t.Errorf("Expected ref value to be 100, got %d", ref.Get())
		}
	})

	t.Run("Test nil Ref", func(t *testing.T) {
		ref := std.NewRef[int](nil)
		result := ref.Get()
		if result != 0 {
			t.Errorf("Expected zero value (0), got %d", result)
		}
	})
}
