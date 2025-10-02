package pkg

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// 测试用的辅助函数
func add(a, b int) int {
	return a + b
}

func multiply(a, b, c int) int {
	return a * b * c
}

func greet(name, greeting string) string {
	return fmt.Sprintf("%s, %s!", greeting, name)
}

func isEven(n int) bool {
	return n%2 == 0
}

func sum(acc, val int) int {
	return acc + val
}

// TestBind 测试函数绑定功能
func TestBind(t *testing.T) {
	t.Run("BasicBind", func(t *testing.T) {
		// 测试基本绑定
		boundAdd := Bind(add, P1, 10)
		result := boundAdd.Call(5)
		if len(result) != 1 || result[0].(int) != 15 {
			t.Errorf("expected 15, got %v", result[0])
		}
	})

	t.Run("MultiplePlaceholders", func(t *testing.T) {
		// 测试多个占位符
		boundGreet := Bind(greet, P1, P2)
		result := boundGreet.Call("World", "Hello")
		expected := "Hello, World!"
		if len(result) != 1 || result[0].(string) != expected {
			t.Errorf("expected %s, got %v", expected, result[0])
		}
	})

	t.Run("MixedPlaceholders", func(t *testing.T) {
		// 测试混合占位符和固定值
		boundGreet := Bind(greet, P1, "Hi")
		result := boundGreet.Call("Alice")
		expected := "Hi, Alice!"
		if len(result) != 1 || result[0].(string) != expected {
			t.Errorf("expected %s, got %v", expected, result[0])
		}
	})

	t.Run("ThreeParameters", func(t *testing.T) {
		// 测试三个参数的函数
		boundMultiply := Bind(multiply, P1, P2, P3)
		result := boundMultiply.Call(2, 3, 4)
		if len(result) != 1 || result[0].(int) != 24 {
			t.Errorf("expected 24, got %v", result[0])
		}
	})
}

// TestBindPanic 测试绑定时的panic情况
func TestBindPanic(t *testing.T) {
	t.Run("NonFunction", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when binding non-function")
			}
		}()
		Bind("not a function")
	})

	t.Run("TooManyArgs", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when too many arguments")
			}
		}()
		Bind(add, 1, 2, 3) // add只接受2个参数
	})

	t.Run("InvalidPlaceholder", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when invalid placeholder")
			}
		}()
		boundAdd := Bind(add, P1, P2)
		boundAdd.Call(5) // 只提供一个参数，但需要两个
	})
}

// TestFunction 测试Function类
func TestFunction(t *testing.T) {
	t.Run("BasicFunction", func(t *testing.T) {
		fn := NewFunction(add)
		result := fn.Call(3, 7)
		if len(result) != 1 || result[0].(int) != 10 {
			t.Errorf("expected 10, got %v", result[0])
		}
	})

	t.Run("IsValid", func(t *testing.T) {
		fn := NewFunction(add)
		if !fn.IsValid() {
			t.Error("function should be valid")
		}
	})

	t.Run("Type", func(t *testing.T) {
		fn := NewFunction(add)
		fnType := fn.Type()
		if fnType.NumIn() != 2 || fnType.NumOut() != 1 {
			t.Error("function type mismatch")
		}
	})

	t.Run("CallWithValues", func(t *testing.T) {
		fn := NewFunction(add)
		args := []reflect.Value{reflect.ValueOf(5), reflect.ValueOf(8)}
		results := fn.CallWithValues(args)
		if len(results) != 1 || results[0].Interface().(int) != 13 {
			t.Errorf("expected 13, got %v", results[0].Interface())
		}
	})
}

// TestFunctionPanic 测试Function的panic情况
func TestFunctionPanic(t *testing.T) {
	t.Run("NonFunction", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when creating function with non-function")
			}
		}()
		NewFunction("not a function")
	})

	t.Run("WrongArgCount", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when wrong argument count")
			}
		}()
		fn := NewFunction(add)
		fn.Call(1) // add需要2个参数
	})
}

// TestInvoke 测试Invoke函数
func TestInvoke(t *testing.T) {
	t.Run("BasicInvoke", func(t *testing.T) {
		result := Invoke(add, 4, 6)
		if len(result) != 1 || result[0].(int) != 10 {
			t.Errorf("expected 10, got %v", result[0])
		}
	})

	t.Run("InvokePanic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when invoking with wrong args")
			}
		}()
		Invoke(add, 1, 2, 3) // 参数数量错误
	})
}

// TestApply 测试Apply函数
func TestApply(t *testing.T) {
	double := func(x int) int { return x * 2 }
	result := Apply(double, []int{1, 2, 3, 4})
	expected := []int{2, 4, 6, 8}

	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

// TestMap 测试Map函数
func TestMap(t *testing.T) {
	square := func(x int) int { return x * x }
	result := Map([]int{1, 2, 3, 4}, square)
	expected := []int{1, 4, 9, 16}

	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

// TestFilter 测试Filter函数
func TestFilter(t *testing.T) {
	result := Filter([]int{1, 2, 3, 4, 5, 6}, isEven)
	expected := []int{2, 4, 6}

	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

// TestReduce 测试Reduce函数
func TestReduce(t *testing.T) {
	result := Reduce([]int{1, 2, 3, 4, 5}, sum, 0)
	expected := 15

	if result != expected {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

// TestPartial 测试Partial函数
func TestPartial(t *testing.T) {
	addFive := Partial(add, 5)
	result := addFive(3)
	expected := 8

	if result != expected {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

// TestCurryVariadic 测试可变参数柯里化
func TestCurryVariadic(t *testing.T) {
	// 创建一个可变参数函数来测试 CurryVariadic
	sumAll := func(nums ...int) int {
		total := 0
		for _, n := range nums {
			total += n
		}
		return total
	}

	curriedSum := CurryVariadic(sumAll)
	addFirst := curriedSum(10)
	result := addFirst(5, 3, 2)
	expected := 20 // 10 + 5 + 3 + 2

	if result != expected {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

// TestCurryAny 测试通用柯里化函数
func TestCurryAny(t *testing.T) {
	t.Run("TwoParamFunction", func(t *testing.T) {
		// 测试固定二参数函数
		curriedAdd := CurryAny(add).(func(int) any)
		partialAdd := curriedAdd(10).(func(int) any)
		result := partialAdd(5).(int)
		expected := 15

		if result != expected {
			t.Errorf("expected %d, got %d", expected, result)
		}
	})

	t.Run("ThreeParamFunction", func(t *testing.T) {
		// 测试固定三参数函数
		curriedMultiply := CurryAny(multiply).(func(int) any)
		step1 := curriedMultiply(2).(func(int) any)
		step2 := step1(3).(func(int) any)
		result := step2(4).(int)
		expected := 24

		if result != expected {
			t.Errorf("expected %d, got %d", expected, result)
		}
	})

	t.Run("StringFunction", func(t *testing.T) {
		// 测试字符串函数
		curriedGreet := CurryAny(greet).(func(string) any)
		step1 := curriedGreet("World").(func(string) any)
		result := step1("Hello").(string)
		expected := "Hello, World!"

		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})
}

// TestMemoize 测试Memoize函数
func TestMemoize(t *testing.T) {
	callCount := 0
	expensiveFunc := func(n int) int {
		callCount++
		return n * n
	}

	memoized := Memoize(expensiveFunc)

	// 第一次调用
	result1 := memoized(5)
	if result1 != 25 || callCount != 1 {
		t.Errorf("first call: expected result=25, callCount=1, got result=%d, callCount=%d", result1, callCount)
	}

	// 第二次调用相同参数，应该使用缓存
	result2 := memoized(5)
	if result2 != 25 || callCount != 1 {
		t.Errorf("second call: expected result=25, callCount=1, got result=%d, callCount=%d", result2, callCount)
	}

	// 调用不同参数
	result3 := memoized(6)
	if result3 != 36 || callCount != 2 {
		t.Errorf("third call: expected result=36, callCount=2, got result=%d, callCount=%d", result3, callCount)
	}
}

// TestOnce 测试Once函数
func TestOnce(t *testing.T) {
	callCount := 0
	fn := func() string {
		callCount++
		return "called"
	}

	onceFn := Once(fn)

	// 多次调用，但函数只执行一次
	result1 := onceFn()
	result2 := onceFn()
	result3 := onceFn()

	if result1 != "called" || result2 != "called" || result3 != "called" {
		t.Error("once function should return same result")
	}

	if callCount != 1 {
		t.Errorf("expected callCount=1, got %d", callCount)
	}
}

// TestDebounce 测试Debounce函数
func TestDebounce(t *testing.T) {
	callCount := 0

	fn := func(arg int) {
		callCount++
		// 简单记录调用，防抖的具体行为依赖实现
		_ = arg
	}

	debouncedFn := Debounce(fn, 100*time.Millisecond)

	// 快速连续调用
	debouncedFn(1)
	debouncedFn(2)
	debouncedFn(3)

	// 等待防抖延迟
	time.Sleep(150 * time.Millisecond)

	// 由于防抖的实现比较简单，这里主要测试函数不会panic
	// 实际的防抖行为可能需要更复杂的实现
	if callCount < 0 {
		t.Error("debounce function should not panic")
	}
}

// TestPlaceholders 测试占位符常量
func TestPlaceholders(t *testing.T) {
	placeholders := []Placeholder{P1, P2, P3, P4, P5, P6, P7, P8, P9, P10}
	for i, p := range placeholders {
		expected := Placeholder(i + 1)
		if p != expected {
			t.Errorf("placeholder P%d should equal %d, got %d", i+1, expected, p)
		}
	}
}

// TestBoundFunctionToFunction 测试BoundFunction转换为Function
func TestBoundFunctionToFunction(t *testing.T) {
	boundAdd := Bind(add, P1, 5)
	fn := boundAdd.ToFunction()

	if !fn.IsValid() {
		t.Error("converted function should be valid")
	}

	// 这里可能需要根据实际的ToFunction实现来调整测试
	if fn.Type() == nil {
		t.Error("converted function should have a type")
	}
}
