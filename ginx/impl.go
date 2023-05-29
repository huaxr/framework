package ginx

import "github.com/gin-gonic/gin"


// register your Router format
type Router interface {
	Router(router *gin.Engine)
}
