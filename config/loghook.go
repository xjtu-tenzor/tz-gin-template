package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"github.com/sirupsen/logrus"
)

type TraceHook struct{}

func (hook *TraceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *TraceHook) Fire(entry *logrus.Entry) error {
	//debug模式升级为 trace模式，方便调试
	if entry.Level == logrus.DebugLevel {
		entry.Level = logrus.TraceLevel
	}

	// 输出堆栈信息：
	stackBuf := make([]byte, 1024)
	stackSize := runtime.Stack(stackBuf, true)
	fmt.Printf("Stack trace:\n%s\n", stackBuf[:stackSize])

	// 输出内存信息
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Println("CPU Num:", runtime.NumCPU())
	fmt.Println("Goroutine Num:", runtime.NumGoroutine())
	fmt.Printf("Memory Stats:\n")
	fmt.Printf("\tHeap Alloc = %v MiB\n", bToMb(memStats.HeapAlloc))
	fmt.Printf("\tHeap Sys = %v MiB\n", bToMb(memStats.HeapSys))
	fmt.Printf("\tHeap Idle = %v MiB\n", bToMb(memStats.HeapIdle))
	fmt.Printf("\tHeap Inuse = %v MiB\n", bToMb(memStats.HeapInuse))
	fmt.Printf("\tHeap Released = %v MiB\n", bToMb(memStats.HeapReleased))
	fmt.Printf("\tHeap Objects = %v\n", memStats.HeapObjects)
	fmt.Printf("\tStack Inuse = %v MiB\n", bToMb(memStats.StackInuse))
	fmt.Printf("\tStack Sys = %v MiB\n", bToMb(memStats.StackSys))
	fmt.Printf("\tMSpan Inuse = %v MiB\n", bToMb(memStats.MSpanInuse))
	fmt.Printf("\tMSpan Sys = %v MiB\n", bToMb(memStats.MSpanSys))
	fmt.Printf("\tMCache Inuse = %v MiB\n", bToMb(memStats.MCacheInuse))
	fmt.Printf("\tMCache Sys = %v MiB\n", bToMb(memStats.MCacheSys))
	fmt.Printf("\tBuckHash Sys = %v MiB\n", bToMb(memStats.BuckHashSys))
	fmt.Printf("\tGCSys = %v MiB\n", bToMb(memStats.GCSys))
	fmt.Printf("\tOther Sys = %v MiB\n", bToMb(memStats.OtherSys))
	fmt.Printf("\tNext GC = %v MiB\n", bToMb(memStats.NextGC))
	fmt.Printf("\tLast GC = %v\n", memStats.LastGC)
	fmt.Printf("\tPause Total Ns = %v\n", memStats.PauseTotalNs)
	fmt.Printf("\tNum GC = %v\n", memStats.NumGC)
	fmt.Printf("\tNum Forced GC = %v\n", memStats.NumForcedGC)
	fmt.Printf("\tGCCPU Fraction = %v\n", memStats.GCCPUFraction)

	SkipSignalChan <- struct{}{}
	return nil
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

type RemoteHook struct {
	Endpoint string
}

func (hook *RemoteHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *RemoteHook) Fire(entry *logrus.Entry) error {
	go func(entry *logrus.Entry) { // 将日志条目转换为 JSON
		logData, err := json.Marshal(entry.Data)
		if err != nil {
			logrus.Errorf("Failed to send log entry to remote server: %v", err)
			return
		}
		// 发送日志到远程服务器
		resp, err := http.Post(hook.Endpoint, "application/json", bytes.NewBuffer(logData))
		if err != nil {
			logrus.Errorf("Failed to send log entry to remote server: %v", err)
			return
		}
		defer resp.Body.Close()
	}(entry)

	return nil
}
