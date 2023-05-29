// Author: huaxinrui@tal.com
// Time: 2022-12-03 17:47
// Git: huaxr

package orm

import (
	"context"

	"github.com/huaxr/framework/internal/metric"

	"xorm.io/xorm/contexts"
)

type defaultHook struct{}

func (*defaultHook) BeforeProcess(c *contexts.ContextHook) (context.Context, error) {
	return c.Ctx, nil
}

func (*defaultHook) AfterProcess(c *contexts.ContextHook) error {
	metric.IncCountWithClear(metric.Mysql)
	return c.Err
}
