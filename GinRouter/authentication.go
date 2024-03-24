package ginrouter

import (
	"net/http"
	ginconfig "sing2cat_web/GinConfig"
	gincontroller "sing2cat_web/GinController"
	middlewarefunc "sing2cat_web/MiddlewareFunc"

	"github.com/gin-gonic/gin"
)

func Authentication_router(r *gin.RouterGroup,html string) {
	user := r.Group("user")
	user.POST("login",func(ctx *gin.Context) {
		secret := ctx.PostForm("secret")
		if !gincontroller.Valid_auth(secret){
			ctx.JSON(http.StatusUnauthorized,gin.H{"error":"unauthorized"})
			return
		}
		token,err := middlewarefunc.Generate_token()
		if err != nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"generate failed"})
			return
		}
		ctx.JSON(http.StatusOK,gin.H{"result":"success","token":token})
	})

	user.GET("verify",middlewarefunc.Jwt_auth(),func(ctx *gin.Context) {
		secret := ctx.GetString("token")
		if !gincontroller.Valid_auth(secret){
			ctx.JSON(http.StatusUnauthorized,gin.H{"error":"unauthorized"})
			return
		}
		token,err := middlewarefunc.Generate_token()
		if err != nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"generate failed"})
			return
		}
		ctx.JSON(http.StatusOK,gin.H{"result":"success","token":token})
	})

	user.POST("/email",func(c *gin.Context) {
		email := c.PostForm("email")
		user_email,_ := ginconfig.Get_value("config","jwt","email")
		if email != user_email{
			c.JSON(http.StatusBadRequest,gin.H{"error":"wrong email"})
			return
		}
		if err := gincontroller.Send_reset_url(email,html);err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusOK,gin.H{"result":"success"})
	})

	user.POST("/reset",func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		captcha := c.PostForm("captcha")
		result,err := gincontroller.Edit_password(email,password,captcha);
		if result && err == nil{
			c.JSON(http.StatusOK,gin.H{"result":"success"})
		}else{
			if err == nil{
				c.JSON(http.StatusInternalServerError,gin.H{"error":"captcha wrong"})
			}else{
				c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			}
		}		
	})
}