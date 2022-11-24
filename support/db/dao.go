package db

import (
	"time"
	"timer/core"
)

func QueryNew(id int) ([]GormDelayTask, error) {
	var list []GormDelayTask
	if err := Manager.Orm.Where("id > ? and status = ?", id, core.Delay).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func QueryByStatus(status core.TaskStatus) ([]GormDelayTask, error) {
	var list []GormDelayTask
	if err := Manager.Orm.Where("status = ?", status).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func QueryTimeout(deadTime int64) ([]GormDelayTask, error) {
	var list []GormDelayTask
	if err := Manager.Orm.Where("status = ? ans update_time < ?", core.Processing, deadTime).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func QueryById(id int) (*GormDelayTask, error) {
	var task GormDelayTask
	if err := Manager.Orm.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateStatus(id int64, status core.TaskStatus, version int) (int64, error) {
	now := time.Now().UnixNano() / 1e6
	result := Manager.Orm.Exec("UPDATE timer_delay_task SET `status` = ? and update_time = ? , version=version+1 WHERE id = ? and version=?",
		status, now, id, version)
	return result.RowsAffected, result.Error
}

func CleanOld(time int64) (int64, error) {
	result := Manager.Orm.Where("update_time < ? and status = ?", time, core.Success).Delete(&GormDelayTask{})
	return result.RowsAffected, result.Error
}

func Insert(task *GormDelayTask) error {
	now := time.Now().UnixNano() / 1e6
	task.CreateTime = now
	task.UpdateTime = now

	if err := Manager.Orm.Create(&task).Error; err != nil {
		return err
	}
	return nil
}
