package pkg

import (
	"fmt"
	"reflect"
)

// std::bind(fn, args...)
// 占位符可以用来表示参数位置, 这样初始化函数时可以指定哪些参数是占位符
// 然后在调用时替换为实际参数(动态绑定)

// 这里使用方法:
// std.Bind(fn, P1, P2, args...) // P1, P2 是占位符,占位符可以放在其他位置
// 绑定函数时, 可以指定占位符, 在调用时替换为实际参数
// 最后使用 Call 方法调用绑定的函数, 并传入实际参数

//例如: f1:= func(a int, b string, c float64) (int, string) { ...}
// f2 := std.Bind(f1, P1, "fixed", P2) // P1, P2 是占位符
// f2.Call(42, 3.14) // 调用时替换占位符, P1 -> 42, P2 -> 3.14

// Placeholder 用于表示绑定时的占位符
type Placeholder int

var (
	_1  = Placeholder(1)
	_2  = Placeholder(2)
	_3  = Placeholder(3)
	_4  = Placeholder(4)
	_5  = Placeholder(5)
	_6  = Placeholder(6)
	_7  = Placeholder(7)
	_8  = Placeholder(8)
	_9  = Placeholder(9)
	_10 = Placeholder(10)

	// 公共占位符变量，用于外部包访问
	P1  = _1
	P2  = _2
	P3  = _3
	P4  = _4
	P5  = _5
	P6  = _6
	P7  = _7
	P8  = _8
	P9  = _9
	P10 = _10
)

// BoundFunction 绑定后的函数类
type BoundFunction struct {
	originalFn   reflect.Value
	boundArgs    []interface{}
	placeholders map[int]int // placeholder -> argument position mapping
	resultType   reflect.Type
}

// std::bind
func Bind(fn interface{}, args ...interface{}) *BoundFunction {
	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		panic("Bind: first argument must be a function")
	}

	fnType := fnVal.Type()
	if len(args) > fnType.NumIn() {
		panic(fmt.Sprintf("Bind: too many arguments, expected at most %d, got %d", fnType.NumIn(), len(args)))
	}

	// 分析占位符
	placeholders := make(map[int]int)
	boundArgs := make([]interface{}, len(args))

	for i, arg := range args {
		if ph, ok := arg.(Placeholder); ok {
			placeholders[int(ph)] = i
		}
		boundArgs[i] = arg
	}

	return &BoundFunction{
		originalFn:   fnVal,
		boundArgs:    boundArgs,
		placeholders: placeholders,
		resultType:   fnType,
	}
}

// Call 调用绑定的函数
func (bf *BoundFunction) Call(args ...interface{}) []interface{} {
	// 准备最终的参数列表
	finalArgs := make([]interface{}, len(bf.boundArgs))
	copy(finalArgs, bf.boundArgs)

	// 替换占位符 - 使用占位符编号对应参数位置
	for i, arg := range finalArgs {
		if ph, ok := arg.(Placeholder); ok {
			phIndex := int(ph) - 1 // 占位符从1开始，参数从0开始
			if phIndex >= len(args) || phIndex < 0 {
				panic(fmt.Sprintf("Bind: placeholder P%d requires argument at index %d, but only %d arguments provided", ph, phIndex+1, len(args)))
			}
			finalArgs[i] = args[phIndex]
		}
	}

	// 转换为 reflect.Value
	in := make([]reflect.Value, len(finalArgs))
	for i, arg := range finalArgs {
		in[i] = reflect.ValueOf(arg)
	}

	// 调用原始函数
	results := bf.originalFn.Call(in)

	// 转换返回值
	out := make([]interface{}, len(results))
	for i, result := range results {
		out[i] = result.Interface()
	}

	return out
}

// CallWithValues 使用 reflect.Value 调用
func (bf *BoundFunction) CallWithValues(args []reflect.Value) []reflect.Value {
	// 准备最终的参数列表
	finalArgs := make([]reflect.Value, len(bf.boundArgs))

	argIndex := 0
	for i, arg := range bf.boundArgs {
		if _, ok := arg.(Placeholder); ok {
			if argIndex >= len(args) {
				panic("Bind: not enough arguments for placeholders")
			}
			finalArgs[i] = args[argIndex]
			argIndex++
		} else {
			finalArgs[i] = reflect.ValueOf(arg)
		}
	}

	return bf.originalFn.Call(finalArgs)
}

// ToFunction 将绑定的函数转换为 Function
func (bf *BoundFunction) ToFunction() *Function[interface{}] {
	return &Function[interface{}]{
		fn:    bf,
		fnVal: reflect.ValueOf(bf.Call),
		fnTyp: reflect.TypeOf(bf.Call),
	}
}
