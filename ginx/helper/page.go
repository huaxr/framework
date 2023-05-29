// Author: huaxr
// Time: 2022-10-24 13:35
// Git: huaxr

package helper

import (
	"strconv"

	"github.com/huaxr/framework/logx"

	"github.com/gin-gonic/gin"
)

func GetPagination(c *gin.Context) (page int, pageSize int) {
	var err error
	page, err = strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		logx.L(c).Errorf("GetPagination page error:%v", err)
		page = 1
	}
	pageSize, err = strconv.Atoi(c.DefaultQuery("page_size", "50"))
	if err != nil || pageSize < 0 {
		logx.L(c).Errorf("GetPagination page_size error:%v", err)
		pageSize = 50
	}
	return
}

// mysql  offset limit generate by gin ctx
func GetOffsetLimit(c *gin.Context) (int, int) {
	page, pageSize := GetPagination(c)
	offset := pageSize * (page - 1)
	return offset, pageSize
}
