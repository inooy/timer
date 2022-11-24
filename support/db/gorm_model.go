package db

// GormDelayTask 延迟任务
type GormDelayTask struct {
	Id         int64  `json:"id" gorm:"primaryKey;comment:消息id"`
	Topic      string `json:"topic" gorm:"column:topic;comment:主题，可以用作业务名称"`
	Body       string `json:"body" gorm:"column:body;comment:消息体"`
	Retry      int    `json:"retry" gorm:"column:retry;comment:失败重试次数，-1表示一直重试，0表示不重试"`
	RetryCount int    `json:"retryCount" gorm:"column:retry_count;comment:已经重试次数"`
	Delay      int64  `json:"delay" gorm:"column:delay;comment:任务延迟时间，单位：毫秒"`
	Status     int    `json:"status" gorm:"column:status;comment:消息状态"`
	CreateTime int64  `json:"createTime" gorm:"column:status;comment:任务创建时间"`
	UpdateTime int64  `json:"update_time" gorm:"column:status;comment:任务更新时间"`
	Version    int    `json:"version" gorm:"column:version;comment:数据版本：乐观锁、防重幂等操作"`
}

// TableName 表名
func (u *GormDelayTask) TableName() string {
	return "timer_delay_task"
}
