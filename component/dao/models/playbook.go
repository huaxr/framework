package models

import (
	"time"
)

type Playbook struct {
	Id          int       `xorm:"not null pk autoincr comment('主键') INT(11)"`
	AppId       int       `xorm:"not null comment('对应的app id') INT(11)"`
	SnapshotId  int       `xorm:"not null comment('快照版本id，用于动态切换') INT(11)"`
	User        string    `xorm:"not null comment('创建人') VARCHAR(255)"`
	Name        string    `xorm:"not null comment('剧本名称') VARCHAR(255)"`
	Enable      int       `xorm:"not null default 0 comment('是否开启') TINYINT(255)"`
	Description string    `xorm:"not null comment('剧本描述') VARCHAR(255)"`
	Token       string    `xorm:"not null comment('授权token，用于被远程调用时鉴权') VARCHAR(255)"`
	UpdateTime  time.Time `xorm:"not null comment('更新时间') DATETIME"`
}
