package main

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type SpinLock int32

func NewSpinLock() sync.Locker {
	var lock SpinLock
	return &lock
}

func (sl *SpinLock) Lock() {
	for !atomic.CompareAndSwapInt32((*int32)(sl), 0, 1) {
		runtime.Gosched()
	}
}

func (sl *SpinLock) Unlock() {
	atomic.StoreInt32((*int32)(sl), 0)
}
