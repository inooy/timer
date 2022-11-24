package core

// DelayTask 延迟任务
type DelayTask struct {
	Id         string     `json:"id"`         // 消息id
	Topic      string     `json:"topic"`      // 主题，可以用作业务名称
	Body       string     `json:"body"`       // 消息体
	Retry      int        `json:"retry"`      // 失败重试次数，-1表示一直重试，0表示不重试
	RetryCount int        `json:"retryCount"` // 已经重试次数
	Delay      int64      `json:"delay"`      //任务延迟时间，单位：毫秒
	Status     TaskStatus `json:"status"`     //消息状态
	CreateTime int64      `json:"createTime"` // 消息创建时间
	Version    int        `json:"version"`    // 数据版本：乐观锁、防重幂等操作
}

type TaskStatus int

const (
	Delay      TaskStatus = 1 // 延迟队列中
	Ready      TaskStatus = 2 //就绪队列中
	Processing TaskStatus = 3 // 消费处理中
	Success    TaskStatus = 4 // 消费成功：保留一定时间，然后批量清理
	Restore    TaskStatus = 5 // 重试队列，失败自动重试队列
	Dead       TaskStatus = 6 // 死亡队列，超过重试次数，需要人工处理
)
