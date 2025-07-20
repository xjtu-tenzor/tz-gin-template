// 提供类似 C++ STL 的功能，包括 function、bind、forward 等
// 我测你的这个代码被AI污染了, 我已经看不懂是什么了
package std

import (
	"reflect"
)

// Version 包版本
const Version = "1.0.0"

// Utility functions and helpers

// Invoke 调用任意可调用对象
func Invoke(fn interface{}, args ...interface{}) []interface{} {
	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		panic("Invoke: argument must be a function")
	}

	// 转换参数
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	// 调用函数
	results := fnVal.Call(in)

	// 转换返回值
	out := make([]interface{}, len(results))
	for i, result := range results {
		out[i] = result.Interface()
	}

	return out
}

// Apply 应用函数到参数列表
func Apply[T any, R any](fn func(T) R, args []T) []R {
	results := make([]R, len(args))
	for i, arg := range args {
		results[i] = fn(arg)
	}
	return results
}

// Map 映射函数（函数式编程）
func Map[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}
	return result
}

// Filter 过滤函数
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Reduce 归约函数
func Reduce[T any, R any](slice []T, fn func(R, T) R, initial R) R {
	result := initial
	for _, item := range slice {
		result = fn(result, item)
	}
	return result
}

// Partial 部分应用函数
func Partial[T1, T2, R any](fn func(T1, T2) R, arg1 T1) func(T2) R {
	return func(arg2 T2) R {
		return fn(arg1, arg2)
	}
}

// Curry 柯里化函数
func Curry2[T1, T2, R any](fn func(T1, T2) R) func(T1) func(T2) R {
	return func(arg1 T1) func(T2) R {
		return func(arg2 T2) R {
			return fn(arg1, arg2)
		}
	}
}

// Curry3 三参数柯里化
func Curry3[T1, T2, T3, R any](fn func(T1, T2, T3) R) func(T1) func(T2) func(T3) R {
	return func(arg1 T1) func(T2) func(T3) R {
		return func(arg2 T2) func(T3) R {
			return func(arg3 T3) R {
				return fn(arg1, arg2, arg3)
			}
		}
	}
}

// Memoize 记忆化函数
func Memoize[T comparable, R any](fn func(T) R) func(T) R {
	cache := make(map[T]R)
	return func(arg T) R {
		if result, ok := cache[arg]; ok {
			return result
		}
		result := fn(arg)
		cache[arg] = result
		return result
	}
}

// Debounce 防抖函数
func Debounce[T any](fn func(T), args T) func() {
	return func() {
		fn(args)
	}
}

// Once 只执行一次的函数
func Once[T any](fn func() T) func() T {
	var (
		called bool
		result T
	)
	return func() T {
		if !called {
			result = fn()
			called = true
		}
		return result
	}
}
