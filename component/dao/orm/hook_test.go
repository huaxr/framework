// Author: huaxr
// Time: 2022-12-07 11:54
// Git: huaxr

package orm

import (
	"testing"

	"github.com/huaxr/framework/component/dao/models"
)

func TestHook(t *testing.T) {
	InitDbInstances()
	orm1, _ := GetEngine()
	var res = make([]models.App, 0)
	err := orm1.Slave().Find(&res)
	t.Log(err, res)

	select {}
}
