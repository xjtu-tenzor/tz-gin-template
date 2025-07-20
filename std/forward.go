package std

import (
	"reflect"
)

// Forward 实现完美转发，类似 C++ 的 std::forward
// 在 Go 中，我们使用接口和反射来实现类似的功能

// Forwarder 转发器接口
type Forwarder interface {
	Forward(args ...interface{}) []interface{}
	ForwardWithValues(args []reflect.Value) []reflect.Value
}

// ForwardingFunction 支持完美转发的函数包装器
type ForwardingFunction struct {
	fn      reflect.Value
	fnType  reflect.Type
	wrapper func([]reflect.Value) []reflect.Value
}

// NewForwardingFunction 创建支持转发的函数
func NewForwardingFunction(fn interface{}) *ForwardingFunction {
	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		panic("NewForwardingFunction: argument must be a function")
	}

	return &ForwardingFunction{
		fn:     fnVal,
		fnType: fnVal.Type(),
	}
}

// Forward 转发参数到目标函数
func (ff *ForwardingFunction) Forward(args ...interface{}) []interface{} {
	// 转换参数为 reflect.Value
	values := make([]reflect.Value, len(args))
	for i, arg := range args {
		values[i] = reflect.ValueOf(arg)
	}

	return ff.ForwardWithValues(values)
}

// ForwardWithValues 使用 reflect.Value 进行转发
func (ff *ForwardingFunction) ForwardWithValues(args []reflect.Value) []interface{} {
	// 调用目标函数
	results := ff.fn.Call(args)

	// 转换返回值
	out := make([]interface{}, len(results))
	for i, result := range results {
		out[i] = result.Interface()
	}

	return out
}

// Chain 链式调用，将多个函数串联起来
func (ff *ForwardingFunction) Chain(next interface{}) *ForwardingFunction {
	// 创建链式调用的函数
	chainedFunc := func(args ...interface{}) interface{} {
		// 调用第一个函数
		firstResults := ff.Forward(args...)

		// 如果第一个函数返回多个值，只取第一个
		var firstResult interface{}
		if len(firstResults) > 0 {
			firstResult = firstResults[0]
		}

		// 调用第二个函数
		nextFn := NewForwardingFunction(next)
		secondResults := nextFn.Forward(firstResult)

		// 返回第二个函数的第一个结果
		if len(secondResults) > 0 {
			return secondResults[0]
		}
		return nil
	}

	return NewForwardingFunction(chainedFunc)
}

// Compose 函数组合，从右到左执行
func Compose(fns ...interface{}) *ForwardingFunction {
	if len(fns) == 0 {
		panic("Compose: at least one function required")
	}

	if len(fns) == 1 {
		return NewForwardingFunction(fns[0])
	}

	// 创建组合函数（从右到左执行）
	composedFunc := func(args ...interface{}) interface{} {
		var result interface{}

		// 最后一个函数接受原始参数
		ff := NewForwardingFunction(fns[len(fns)-1])
		results := ff.Forward(args...)
		if len(results) > 0 {
			result = results[0]
		}

		// 从右到左逐个执行
		for i := len(fns) - 2; i >= 0; i-- {
			ff := NewForwardingFunction(fns[i])
			results := ff.Forward(result)
			if len(results) > 0 {
				result = results[0]
			}
		}

		return result
	}

	return NewForwardingFunction(composedFunc)
}

// Pipe 管道操作，从左到右执行
func Pipe(fns ...interface{}) *ForwardingFunction {
	if len(fns) == 0 {
		panic("Pipe: at least one function required")
	}

	if len(fns) == 1 {
		return NewForwardingFunction(fns[0])
	}

	// 创建管道函数
	pipelineFunc := func(args ...interface{}) interface{} {
		var result interface{}

		// 第一个函数接受原始参数
		ff := NewForwardingFunction(fns[0])
		results := ff.Forward(args...)
		if len(results) > 0 {
			result = results[0]
		}

		// 后续函数逐个处理
		for i := 1; i < len(fns); i++ {
			ff := NewForwardingFunction(fns[i])
			results := ff.Forward(result)
			if len(results) > 0 {
				result = results[0]
			}
		}

		return result
	}

	return NewForwardingFunction(pipelineFunc)
}

// Move 实现移动语义（在 Go 中主要用于大对象的传递优化）
func Move[T any](value T) T {
	// 在 Go 中，我们不能真正实现 C++ 的移动语义
	// 但可以通过一些约定来优化性能
	return value
}

// Forward 通用转发函数
func Forward[T any](value T) T {
	// 在 Go 中，参数传递本身就是值传递或引用传递
	// 这个函数主要用于语义上的完美转发
	return value
}

// Ref 创建引用包装器
type Ref[T any] struct {
	Value *T
}

// NewRef 创建新的引用
func NewRef[T any](value *T) Ref[T] {
	return Ref[T]{Value: value}
}

// Get 获取引用的值
func (r Ref[T]) Get() T {
	if r.Value == nil {
		var zero T
		return zero
	}
	return *r.Value
}

// Set 设置引用的值
func (r Ref[T]) Set(value T) {
	if r.Value != nil {
		*r.Value = value
	}
}
