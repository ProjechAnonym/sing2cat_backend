package ginrouter

import (
	"net/http"
	gincontroller "sing2cat_web/GinController"
	middlewarefunc "sing2cat_web/MiddlewareFunc"

	"github.com/gin-gonic/gin"
)

func Fetch_router(r *gin.RouterGroup) {
	fetch_router := r.Group("fetch")
	fetch_router.Use(middlewarefunc.Jwt_auth())
	// 获取网站组信息
	fetch_router.GET("/component", func(ctx *gin.Context) {
		components, err := gincontroller.Fetch_components()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": components})
	})

}