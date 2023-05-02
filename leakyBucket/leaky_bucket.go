package main

import (
	"math"
	"sync"
	"time"
)

type LeakyBucket struct {
	rate     float64 // 每秒固定出水速率
	capacity float64 // 桶的容量
	water    float64 // 当前水量

	lastLeakMs int64 // 桶上次漏水时间戳 ms
	lock       sync.Mutex
}

func (l *LeakyBucket) Pass() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now().UnixNano() / 1e6
	e := float64(now-l.lastLeakMs) * l.rate / 1000 // 先执行漏水
	l.water = l.water - e                          // 计算剩余水量
	l.water = math.Max(0, l.water)                 // 桶空了
	l.lastLeakMs = now
	if (l.water + 1) < l.capacity {
		// 加水
		l.water++
		return true
	} else {
		// 水满，拒绝加水
		return false
	}
}

func (l *LeakyBucket) Set(r, c float64) {
	l.rate = r
	l.capacity = c
	l.water = 0
	l.lastLeakMs = time.Now().UnixNano() / 1e6
}
