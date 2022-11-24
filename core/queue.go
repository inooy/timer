package core

type DelayQueue interface {
	GetName() string
	Push(task *DelayTask) error
	ReToDelayQueue(task *DelayTask) error
}
