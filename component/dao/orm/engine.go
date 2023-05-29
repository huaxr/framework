// Author: XinRui Hua
// Time:   2022/4/12 下午2:18
// Git:    huaxr

package orm

import (
	"context"
	"fmt"
	"time"

	"github.com/huaxr/framework/logx"

	"github.com/huaxr/framework/pkg/confutil"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var (
	egs = make([]*xorm.EngineGroup, 0)
)

func InitDbInstances() error {
	for _, mysql := range confutil.GetDefaultConfig().Mysql {
		i := mysql
		var err error
		master, err := xorm.NewEngine("mysql", i.Master)
		if err != nil {
			return err
		}

		slaves := make([]*xorm.Engine, 0)
		for _, slave := range i.Slaves {
			s, err := xorm.NewEngine("mysql", slave)
			if err != nil {
				return err
			}
			slaves = append(slaves, s)
		}

		eg, err := xorm.NewEngineGroup(master, slaves)
		if err != nil {
			return err
		}

		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		session := eg.Context(ctx)
		if err = session.Ping(); err != nil {
			logx.L().Errorf("err ping mysql %v", err)
			return err
		}

		eg.AddHook(&defaultHook{})
		eg.SetMaxIdleConns(i.MaxIdle)
		eg.SetMaxOpenConns(i.MaxConn)
		eg.SetLogger(&dbLog{
			slowDuration: time.Duration(i.SlowDuration) * time.Millisecond,
			showSql:      i.ShowSql,
		})

		egs = append(egs, eg)
	}

	return nil
}

func GetEngineByIndex(index int) (*xorm.EngineGroup, error) {
	if len(egs)-1 < index {
		return nil, fmt.Errorf("index out of range")
	}
	return egs[index], nil
}

func GetEngine() (*xorm.EngineGroup, error) {
	if len(egs) == 0 {
		return nil, fmt.Errorf("no available mysql client")
	}
	return egs[0], nil
}
