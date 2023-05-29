// Author: huaxinrui@tal.com
// Time: 2022-11-10 20:19
// Git: huaxr

package transport

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCtx(t *testing.T) {
	c := &gin.Context{
		Keys: map[string]interface{}{},
	}

	grcpC := CtxToGRpcCtx(c)
	grcpC = CtxToGRpcCtxWithAllField(c)
	grcpC = CtxToGRpcCtxWithManyFields(c, map[string]string{"a": "b"})
	grcpC = CtxToGRpcCtxWithOneField(c, "a", "b")
	grcpC = CtxToGRpcCtxWithField(c, "a", []byte{123})

	t.Fatal(grcpC)

}
