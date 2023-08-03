package core

import (
	"errors"
	"sync"
	"time"
)

type FourHeap struct {
	size      int32
	capacity  int32
	timers    []*Timer
	isStoped  bool
	isStarted bool
	mu        sync.Mutex
}

func NewFourHeap(capacity int32) *FourHeap {
	if capacity <= 0 {
		return nil
	}
	return &FourHeap{
		size:     0,
		capacity: capacity,
		timers:   make([]*Timer, capacity),
	}
}

func (that *FourHeap) Start() {
	that.mu.Lock()
	if that.isStarted {
		return
	}
	that.isStarted = true
	that.mu.Unlock()

	go timerCheck(that)
}

func (that *FourHeap) AddTimer(timer *Timer) error {
	if that.isStoped {
		return errors.New("stop")
	}
	if timer == nil {
		return errors.New("timer is nil")
	}
	if that.size == that.capacity {
		return errors.New("enough")
	}

	that.mu.Lock()
	defer that.mu.Unlock()

	that.timers[that.size] = timer
	that.size += 1
	that.heapUp(that.size - 1)
	return nil
}

func (that *FourHeap) Peek() *Timer {
	if that.size == 0 {
		return nil
	}

	that.mu.Lock()
	defer that.mu.Unlock()

	if that.size == 0 {
		return nil
	}

	return that.timers[0]
}

func (that *FourHeap) Pop() *Timer {
	if that.size == 0 {
		return nil
	}

	that.mu.Lock()
	defer that.mu.Unlock()

	if that.size == 0 {
		return nil
	}
	ans := that.timers[0]
	that.deleteFromHeap(0)
	return ans
}

func (that *FourHeap) deleteFromHeap(targetIndex int32) {
	if targetIndex < 0 || targetIndex == that.size {
		return
	}

	// 交换
	temp := that.timers[targetIndex]
	that.timers[targetIndex] = that.timers[that.size-1]
	that.timers[that.size-1] = temp
	that.size -= 1
	that.heapify(targetIndex)
}

func (that *FourHeap) heapUp(targetIndex int32) {
	if targetIndex < 0 || targetIndex == that.size {
		return
	}

	parentIndex := int32(1)
	for targetIndex != 0 {
		parentIndex = (targetIndex - 1) / 4
		if that.timers[targetIndex].expect < that.timers[parentIndex].expect {
			temp := that.timers[targetIndex]
			that.timers[targetIndex] = that.timers[parentIndex]
			that.timers[parentIndex] = temp
			targetIndex = parentIndex
		} else {
			break
		}
	}
}

func (that *FourHeap) heapify(targetIndex int32) {
	if targetIndex < 0 || targetIndex == that.size {
		return
	}

	leftChildIndex := (targetIndex << 2) + 1
	childMinIndex := int32(1)
	for leftChildIndex < that.size {
		childMinIndex = int32(leftChildIndex)
		for i := leftChildIndex + 1; i < that.size; i++ {
			if that.timers[childMinIndex].expect > that.timers[i].expect {
				childMinIndex = i
			}
		}
		if that.timers[childMinIndex].expect >= that.timers[targetIndex].expect {
			break
		}
		temp := that.timers[childMinIndex]
		that.timers[childMinIndex] = that.timers[targetIndex]
		that.timers[targetIndex] = temp
		targetIndex = childMinIndex
		leftChildIndex = (targetIndex << 2) + 1
	}
}

func timerCheck(h *FourHeap) {
	nowTimestamp := int64(1)
	for !h.isStoped {
		nowTimestamp = time.Now().UnixMilli()
		timerCheckRun(h, nowTimestamp)
	}
}

func timerCheckRun(h *FourHeap, nowTimestamp int64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.size == 0 {
		return
	}

	task := h.timers[0]
	if nowTimestamp < task.expect {
		return
	}
	ans := h.timers[0]
	h.deleteFromHeap(0)
	ans.result = ans.taskFunc(ans.args)
	ans.status = RUN_END
}
