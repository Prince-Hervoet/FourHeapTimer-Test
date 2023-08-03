package core

import "time"

const (
	WAITING = "waiting"
	RUNNING = "running"
	CANCEL  = "cancel"
	ERROR   = "error"
	RUN_END = "end"
)

type Timer struct {
	taskFunc func(any) any
	args     any
	result   any
	expect   int64
	status   string
}

func NewTimer(taskFunc func(any) any, args any, will int64) *Timer {
	return &Timer{
		taskFunc: taskFunc,
		args:     args,
		expect:   time.Now().UnixMilli() + will,
		status:   WAITING,
	}
}

func (that *Timer) GetExpect() int64 {
	return that.expect
}

func (that *Timer) GetStatus() string {
	return that.status
}

func (that *Timer) GetResult() any {
	if that.status != RUN_END {
		return nil
	}
	return that.result
}
