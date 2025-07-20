package test

import (
	"template/std"
	"testing"
)

// BenchmarkFunction 测试 Function 的性能
func BenchmarkFunction(b *testing.B) {
	add := func(a, b int) int {
		return a + b
	}

	fn := std.NewFunction(add)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn.Call(i, i+1)
	}
}

// BenchmarkDirectCall 测试直接调用的性能（作为对比）
func BenchmarkDirectCall(b *testing.B) {
	add := func(a, b int) int {
		return a + b
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		add(i, i+1)
	}
}

// BenchmarkBind 测试 Bind 的性能
func BenchmarkBind(b *testing.B) {
	multiply := func(a, b, c int) int {
		return a * b * c
	}

	bound := std.Bind(multiply, 2, std.P1, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bound.Call(i)
	}
}

// BenchmarkForward 测试 Forward 的性能
func BenchmarkForward(b *testing.B) {
	double := func(x int) int {
		return x * 2
	}

	ff := std.NewForwardingFunction(double)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ff.Forward(i)
	}
}

// BenchmarkPipe 测试 Pipe 的性能
func BenchmarkPipe(b *testing.B) {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline.Forward(i)
	}
}

// BenchmarkMap 测试 Map 的性能
func BenchmarkMap(b *testing.B) {
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	double := func(x int) int {
		return x * 2
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		std.Map(numbers, double)
	}
}

// BenchmarkFilter 测试 Filter 的性能
func BenchmarkFilter(b *testing.B) {
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	isEven := func(x int) bool {
		return x%2 == 0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		std.Filter(numbers, isEven)
	}
}

// BenchmarkReduce 测试 Reduce 的性能
func BenchmarkReduce(b *testing.B) {
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	add := func(acc, x int) int {
		return acc + x
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		std.Reduce(numbers, add, 0)
	}
}

// BenchmarkMemoize 测试 Memoize 的性能
func BenchmarkMemoize(b *testing.B) {
	expensiveFunc := func(n int) int {
		// 模拟一些计算
		result := 0
		for i := 0; i <= n%100; i++ { // 限制计算量
			result += i
		}
		return result
	}

	memoized := std.Memoize(expensiveFunc)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memoized(i % 50) // 重复调用以测试缓存效果
	}
}

// BenchmarkMemoizeWithoutCache 测试没有缓存的性能（作为对比）
func BenchmarkMemoizeWithoutCache(b *testing.B) {
	expensiveFunc := func(n int) int {
		// 模拟一些计算
		result := 0
		for i := 0; i <= n%100; i++ {
			result += i
		}
		return result
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expensiveFunc(i % 50)
	}
}

// BenchmarkCurry 测试 Curry 的性能
func BenchmarkCurry(b *testing.B) {
	add := func(a, b int) int {
		return a + b
	}

	curriedAdd := std.Curry2(add)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addFive := curriedAdd(5)
		addFive(i)
	}
}

// BenchmarkOnce 测试 Once 的性能
func BenchmarkOnce(b *testing.B) {
	counter := 0
	increment := func() int {
		counter++
		return counter
	}

	once := std.Once(increment)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		once()
	}
}
