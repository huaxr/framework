package models

import (
	"time"
)

type User struct {
	Id           int       `xorm:"not null pk autoincr comment('主键id') INT(11)"`
	Account      string    `xorm:"not null comment('用户') VARCHAR(255)"`
	Name         string    `xorm:"not null comment('邮箱') VARCHAR(255)"`
	Workcode     string    `xorm:"not null comment('工号') VARCHAR(255)"`
	DeptId       string    `xorm:"not null comment('部门') VARCHAR(255)"`
	DeptName     string    `xorm:"not null comment('部门') VARCHAR(255)"`
	DeptFullName string    `xorm:"not null comment('部门') VARCHAR(255)"`
	Email        string    `xorm:"not null comment('email') VARCHAR(255)"`
	Avatar       string    `xorm:"not null comment('头像') VARCHAR(255)"`
	CreateTime   time.Time `xorm:"not null comment('创建时间') DATETIME"`
	SuperAdmin   int       `xorm:"not null comment('超级管理员') TINYINT(1)"`
}
