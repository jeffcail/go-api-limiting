package main

import (
	"log"
	"sync"
	"time"
)

type Counter struct {
	max          int           // 时间窗口内最大请求数
	firstReqTime time.Time     // 请求开始时间
	tt           time.Duration // 时间窗口
	count        int           // 时间窗口内累计的请求次数
	lock         sync.Mutex
}

func (c *Counter) Pass() bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.count == c.max-1 {
		now := time.Now()
		if now.Sub(c.firstReqTime) >= c.tt {
			c.Reset(now)
			return true
		} else {
			return false
		}
	} else {
		c.count++
		return true
	}
}

func (c *Counter) Reset(t time.Time) {
	c.firstReqTime = t
	c.count = 0
}

func (c *Counter) Set(r int, tt time.Duration) {
	c.max = r
	c.firstReqTime = time.Now()
	c.tt = tt
	c.count = 0
}

func main() {
	var wg sync.WaitGroup
	var lc Counter
	lc.Set(10, time.Second) // 1s内10速率
	for i := 0; i < 50; i++ {
		wg.Add(1)
		log.Println("创建请求: ", i)
		go func(i int) {

			if lc.Pass() {
				log.Println("response: ", i)
			}
			wg.Done()
		}(i)
		time.Sleep(10 * time.Millisecond)
	}
	wg.Wait()
}
