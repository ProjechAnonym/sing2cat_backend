package ginrouter

import (
	"net/http"
	gincontroller "sing2cat_web/GinController"
	middlewarefunc "sing2cat_web/MiddlewareFunc"

	"github.com/gin-gonic/gin"
)

func Remove_router(r *gin.RouterGroup) {
	delete_router := r.Group("delete")
	delete_router.Use(middlewarefunc.Jwt_auth())
	delete_router.POST("/component", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		app := ctx.PostForm("app")
		if err := gincontroller.Delete_component(name,app); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"result": "success"})
	})
}