package std

import (
	"fmt"
	"reflect"
)

// Function 类似于 C++ 的 std::function，用于存储可调用对象
type Function[T any] struct {
	fn    interface{}
	fnVal reflect.Value
	fnTyp reflect.Type
}

// NewFunction 创建一个新的 Function 实例
func NewFunction[T any](fn T) *Function[T] {
	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		panic("NewFunction: argument must be a function")
	}

	return &Function[T]{
		fn:    fn,
		fnVal: fnVal,
		fnTyp: fnVal.Type(),
	}
}

// Call 调用存储的函数
func (f *Function[T]) Call(args ...interface{}) []interface{} {
	if f.fn == nil {
		panic("Function: cannot call nil function")
	}

	// 检查参数数量
	if f.fnTyp.NumIn() != len(args) {
		panic(fmt.Sprintf("Function: expected %d arguments, got %d", f.fnTyp.NumIn(), len(args)))
	}

	// 转换参数为 reflect.Value
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	// 调用函数
	results := f.fnVal.Call(in)

	// 转换返回值
	out := make([]interface{}, len(results))
	for i, result := range results {
		out[i] = result.Interface()
	}

	return out
}

// CallWithValues 使用 reflect.Value 调用函数
func (f *Function[T]) CallWithValues(args []reflect.Value) []reflect.Value {
	if f.fn == nil {
		panic("Function: cannot call nil function")
	}
	return f.fnVal.Call(args)
}

// IsValid 检查函数是否有效
func (f *Function[T]) IsValid() bool {
	return f.fn != nil && f.fnVal.IsValid()
}

// Type 返回函数类型
func (f *Function[T]) Type() reflect.Type {
	return f.fnTyp
}

// 便利函数：创建常见类型的 Function

// Func0 无参数函数
type Func0[R any] func() R

func NewFunc0[R any](fn func() R) *Function[Func0[R]] {
	return NewFunction(Func0[R](fn))
}

// Func1 一个参数的函数
type Func1[T, R any] func(T) R

func NewFunc1[T, R any](fn func(T) R) *Function[Func1[T, R]] {
	return NewFunction(Func1[T, R](fn))
}

// Func2 两个参数的函数
type Func2[T1, T2, R any] func(T1, T2) R

func NewFunc2[T1, T2, R any](fn func(T1, T2) R) *Function[Func2[T1, T2, R]] {
	return NewFunction(Func2[T1, T2, R](fn))
}

// Func3 三个参数的函数
type Func3[T1, T2, T3, R any] func(T1, T2, T3) R

func NewFunc3[T1, T2, T3, R any](fn func(T1, T2, T3) R) *Function[Func3[T1, T2, T3, R]] {
	return NewFunction(Func3[T1, T2, T3, R](fn))
}
