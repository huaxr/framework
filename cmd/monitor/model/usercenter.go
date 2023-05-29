package model

import "time"

type Student struct {
	Id        int32     `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	UserId    int32     `xorm:"int(11)" json:"user_id"`
	StuId     int32     `xorm:"int(11)" json:"stu_id"`
	TalId     string    `xorm:"VARCHAR(100)" json:"tal_id"`
	OrgId     int32     `xorm:"int(11)" json:"org_id"`
	Sn        string    `xorm:"VARCHAR(100)" json:"sn"`
	GradeId   int32     `xorm:"int(11)" json:"grade_id"`
	CreatedAt time.Time `xorm:"created_at created" json:"-"` // 创建时间
	UpdatedAt time.Time `xorm:"updated_at updated" json:"-"` // 修改时间
	DeletedAt time.Time `xorm:"deleted_at deleted" json:"-" `
}
