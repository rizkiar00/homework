package entity

import "github.com/rizkiar00/homework/pkg/constant"

type Action struct {
	ActionID   int64  `gorm:"column:action_id;primaryKey;autoIncrement"`
	ActionDesc string `gorm:"column:action_desc"`
	ActionType string `gorm:"column:action_type"`
	Endpoint   string `gorm:"column:endpoint"`
}

func (Action) TableName() string {
	return constant.TableActions
}
