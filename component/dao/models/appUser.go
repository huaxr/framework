package models

import (
	"time"
)

type AppUser struct {
	Id         int       `xorm:"not null pk autoincr comment('主键') INT(11)"`
	UserId     int       `xorm:"not null comment('外键') INT(11)"`
	AppId      int       `xorm:"not null comment('app外键') INT(11)"`
	Checked    int       `xorm:"not null comment('是否通过') TINYINT(1)"`
	CreateTime time.Time `xorm:"not null comment('创建时间') DATETIME"`
}
