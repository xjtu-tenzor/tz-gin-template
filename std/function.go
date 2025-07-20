package std

import (
	"fmt"
	"reflect"
)

// Function 类似于 C++ 的 std::function，用于存储可调用对象
// 仿函数, 所有可以通过 f_name(args...) 调用的对象, 可以用forward.go中实现
// ↑见func_prog.go
// std::function<return_type(arg1_type, arg2_type, ...)> f_name;

// 有人问这样做不是更麻烦, 但是包装了一层之后函数什么时候调用, 用什么参数调用
// 是否需要先空几个参数, 到时候再放进去

type Function[T any] struct {
	fn    interface{}
	fnVal reflect.Value
	fnTyp reflect.Type
}

// NewFunction
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

// Call 调用存储的函数, 区分一下和反射库的call
// 反射库的call参数类型是 []reflect.Value, auto&& 类型
// 这个封装了一下,直接塞数值就行
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
