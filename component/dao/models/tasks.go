package models

import (
	"time"
)

type Tasks struct {
	Id            int       `xorm:"not null pk autoincr comment('主键id') INT(11)"`
	Name          string    `xorm:"not null default '' comment('函数名') VARCHAR(255)"`
	Configuration string    `xorm:"not null default '' comment('函数的配置') VARCHAR(255)"`
	Xrn           string    `xorm:"not null default '' comment('对应的xrn') VARCHAR(255)"`
	Description   string    `xorm:"not null default '' comment('描述') VARCHAR(255)"`
	AppId         int       `xorm:"not null comment('对应的appid') index INT(11)"`
	Type          string    `xorm:"not null comment('任务类型') VARCHAR(255)"`
	UpdateTime    time.Time `xorm:"not null comment('更新时间') DATETIME"`
	InputExample  string    `xorm:"not null default '' comment('输入样例') VARCHAR(255)"`
	OutputExample string    `xorm:"not null default '' comment('输出样例') VARCHAR(255)"`
	User          string    `xorm:"not null comment('创建者') VARCHAR(255)"`
}
