package models

import (
	"time"
)

type Snapshot struct {
	Id         int       `xorm:"not null pk autoincr comment('主键') INT(11)"`
	PlaybookId int       `xorm:"not null comment('对应的剧本id') index INT(11)"`
	Snapshot   string    `xorm:"comment('剧本快照body体') TEXT"`
	Rawbody    string    `xorm:"comment('前端传来的元数据') TEXT"`
	Checksum   string    `xorm:"not null comment('校验和') index VARCHAR(255)"`
	AppId      int       `xorm:"not null comment('app') index INT(11)"`
	Snapname   string    `xorm:"not null comment('快照名称') VARCHAR(255)"`
	UpdateTime time.Time `xorm:"not null comment('更新时间') DATETIME"`
	User       string    `xorm:"not null comment('更新快照人') VARCHAR(255)"`
}
