package models

import (
	"time"
)

type App struct {
	Id            int       `xorm:"not null pk autoincr comment('主键id') INT(11)"`
	AppName       string    `xorm:"not null default '' comment('app名称') unique VARCHAR(255)"`
	User          string    `xorm:"not null default '0' comment('用户') VARCHAR(255)"`
	Token         string    `xorm:"not null comment('是否有权限消费nsq') unique VARCHAR(255)"`
	Brokers       string    `xorm:"not null comment('brokers') VARCHAR(255)"`
	BrokerType    string    `xorm:"not null comment('消息队列类型') VARCHAR(255)"`
	Eps           int       `xorm:"not null comment('限流') INT(255)"`
	GroupId       string    `xorm:"not null comment('用户字段组') VARCHAR(11)"`
	UpdateTime    time.Time `xorm:"not null comment('更新时间') DATETIME"`
	Checked       int       `xorm:"not null default 0 comment('是否通过核审') TINYINT(1)"`
	Description   string    `xorm:"not null comment('信息描述') VARCHAR(255)"`
	LastAliveTime time.Time `xorm:"not null comment('上次心跳时间') DATETIME"`
	Share         int       `xorm:"not null default 0 comment('共享app') TINYINT(1)"`
}
