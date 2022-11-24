package timer

import (
	"gorm.io/gorm"
	"timer/core"
	"timer/support/db"
)

func SetUp(orm *gorm.DB) *core.DelayQueueManager {
	core.Manager = &core.DelayQueueManager{
		QueueMap: make(map[string]core.DelayQueue),
	}
	db.Setup(orm)
	return core.Manager
}
