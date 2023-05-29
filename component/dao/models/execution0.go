package models

type Execution0 struct {
	Id         int    `xorm:"not null pk autoincr comment('主键') INT(11)"`
	TraceId    string `xorm:"not null comment('执行链路id') index VARCHAR(255)"`
	Sequence   int    `xorm:"not null comment('节点code') INT(11)"`
	NodeCode   string `xorm:"not null comment('节点id') index VARCHAR(32)"`
	Domain     string `xorm:"not null comment('所属域') VARCHAR(255)"`
	Status     string `xorm:"not null comment('状态') VARCHAR(30)"`
	PlaybookId int    `xorm:"not null comment('剧本id') INT(11)"`
	SnapshotId int    `xorm:"not null comment('快照id') index INT(11)"`
	Extra      string `xorm:"comment('额外信息') TEXT"`
	Timestamp  string `xorm:"not null comment('时间戳') VARCHAR(55)"`
	Chain      string `xorm:"not null comment('执行链路') VARCHAR(255)"`
}
