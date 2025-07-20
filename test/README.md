# 测试文件组织结构

这个目录包含了 STD 包的所有测试和示例文件。

## 文件结构

```
test/
├── std_test.go         # 单元测试文件
├── benchmark_test.go   # 基准测试文件
└── examples.go         # 详细示例文件

```

## 运行测试

### 1. 运行所有示例
```bash
go run test_runner.go
```

### 2. 运行单元测试
```bash
# 运行所有测试
go test ./test -v

# 运行特定测试
go test ./test -run TestFunction -v
go test ./test -run TestBind -v
go test ./test -run TestForward -v
```

### 3. 运行基准测试
```bash
# 运行所有基准测试
go test ./test -bench=.

# 运行特定基准测试
go test ./test -bench=BenchmarkFunction
go test ./test -bench=BenchmarkBind

# 运行基准测试并显示内存分配
go test ./test -bench=. -benchmem
```

### 4. 运行 Placeholder 演示
```bash
go run placeholder_demo.go
```

## 测试覆盖率

查看测试覆盖率：
```bash
go test ./test -cover

# 生成详细的覆盖率报告
go test ./test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 测试文件说明

### std_test.go
包含所有功能的单元测试：
- TestFunction: 测试 std::function 功能
- TestBind: 测试 std::bind 功能
- TestForward: 测试 std::forward 功能
- TestFunctionalProgramming: 测试函数式编程特性
- TestAdvancedFeatures: 测试高级功能（记忆化、Once等）
- TestRef: 测试引用功能

### benchmark_test.go
包含性能基准测试：
- BenchmarkFunction: Function 的性能测试
- BenchmarkBind: Bind 的性能测试
- BenchmarkForward: Forward 的性能测试
- BenchmarkMap/Filter/Reduce: 函数式编程的性能测试
- BenchmarkMemoize: 记忆化的性能测试

### examples.go
包含详细的使用示例：
- 基本功能示例
- 复杂使用场景
- 实际应用案例

## 持续集成

如果要在 CI/CD 中运行测试，可以使用以下命令：

```bash
# 运行所有测试（包括基准测试）
go test ./test -v -bench=. -benchmem

# 仅运行单元测试
go test ./test -v

# 运行测试并生成覆盖率报告
go test ./test -v -cover -coverprofile=coverage.out
```

## 添加新测试

当添加新功能时，请遵循以下规则：

1. 在 `std_test.go` 中添加相应的单元测试
2. 在 `benchmark_test.go` 中添加性能测试（如果需要）
3. 在 `examples.go` 中添加使用示例
4. 确保所有测试都能通过：`go test ./test -v`
