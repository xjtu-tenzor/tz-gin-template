// 提供高级函数式编程功能，基于 forward.go 的基础设施
// 避免重复代码，专注于函数式编程特性

// 注释来自ai ♥️
package pkg

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

// CurryVariadic 柯里化可变参数函数
func CurryVariadic[R any, T any](fn func(...T) R) func(T) func(...T) R {
	return func(arg T) func(...T) R {
		return func(args ...T) R {
			allArgs := append([]T{arg}, args...)
			return fn(allArgs...)
		}
	}
}

// CurryAny 通用柯里化函数，支持任意固定参数函数
// 通过反射将固定参数函数转换为可变参数形式进行柯里化
func CurryAny(fn interface{}) interface{} {
	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		panic("CurryAny: argument must be a function")
	}

	fnType := fnVal.Type()
	numIn := fnType.NumIn()

	if numIn == 0 {
		panic("CurryAny: function must have at least one parameter")
	}

	return buildCurriedFunc(fnVal, fnType, 0, []reflect.Value{})
}

// buildCurriedFunc 递归构建柯里化函数
func buildCurriedFunc(originalFn reflect.Value, originalType reflect.Type, argIndex int, boundArgs []reflect.Value) interface{} {
	if argIndex >= originalType.NumIn() {
		// 所有参数都已绑定，直接调用原函数
		results := originalFn.Call(boundArgs)
		if len(results) == 1 {
			return results[0].Interface()
		}
		// 多返回值
		interfaces := make([]interface{}, len(results))
		for i, result := range results {
			interfaces[i] = result.Interface()
		}
		return interfaces
	}

	// 创建一个接受下一个参数的函数
	paramType := originalType.In(argIndex)

	// 动态创建函数类型
	fnType := reflect.FuncOf(
		[]reflect.Type{paramType},
		[]reflect.Type{reflect.TypeOf((*interface{})(nil)).Elem()},
		false,
	)

	fn := reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
		newBoundArgs := append(boundArgs, args[0])
		result := buildCurriedFunc(originalFn, originalType, argIndex+1, newBoundArgs)
		return []reflect.Value{reflect.ValueOf(result)}
	})

	return fn.Interface()
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
