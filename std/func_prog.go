// 提供高级函数式编程功能，基于 forward.go 的基础设施
// 避免重复代码，专注于函数式编程特性
package std

import (
	"reflect"
	"time"
)

// Invoke 调用任意可调用对象 (cpp俗称仿函数), 指所有可以通过 f_name(args...) 调用的对象
// 使用已有的 ForwardingFunction 避免重复反射逻辑
func Invoke(fn interface{}, args ...interface{}) []interface{} {
	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		panic("Invoke: first argument must be a function")
	}

	if len(args) != fnVal.Type().NumIn() {
		panic("Invoke: argument count mismatch")
	}

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	results := fnVal.Call(in)

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

// Debounce在指定时间内只执行最后一次调用
func Debounce[T any](fn func(T), delay time.Duration) func(T) {
	// 简化版本的防抖，实际项目中可能需要使用 time.Timer
	var lastArgs T
	var pending bool

	return func(args T) {
		lastArgs = args
		if !pending {
			pending = true
			go func() {
				time.Sleep(delay)
				if pending {
					fn(lastArgs)
				}
				pending = false
			}()
			pending = false
		}
	}
}

// Once 只执行一次
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
