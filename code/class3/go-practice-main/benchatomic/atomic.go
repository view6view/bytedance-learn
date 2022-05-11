package benchatomic

import (
	"sync"
	"sync/atomic"
)

type atomicCounter struct {
	i int32
}

func AtomicAddOne(c *atomicCounter) {
	atomic.AddInt32(&c.i, 1)
}

type mutexCounter struct {
	i int32
	m sync.Mutex
}

func MutexAddOne(c *mutexCounter) {
	c.m.Lock()
	c.i++
	c.m.Unlock()
}
