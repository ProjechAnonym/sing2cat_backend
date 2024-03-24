package ginrouter

import (
	"encoding/json"
	"net/http"
	ginconfig "sing2cat_web/GinConfig"
	gincontroller "sing2cat_web/GinController"
	middlewarefunc "sing2cat_web/MiddlewareFunc"

	"github.com/gin-gonic/gin"
)

func Add_item(r *gin.RouterGroup) {
	add_router := r.Group("add")
	add_router.Use(middlewarefunc.Jwt_auth())
	add_router.POST("/component",func(ctx *gin.Context) {
		
		var content ginconfig.Component
		// 解析前端的字符串
		if err := ctx.BindJSON(&content);err!=nil{
			ginconfig.Logger_caller("Marshal json failed!",err)
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		// 补充所需数据
		data,_ := json.Marshal(content.Data)
		content.Gorm_data = string(data)
		if err := ginconfig.Db.Create(&content).Error;err!=nil{
			ginconfig.Logger_caller("Write msg to the database failed!",err)
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		ctx.JSON(http.StatusOK,gin.H{"result":"success"})
	})
	// 提交图片路由
	add_router.POST("/pic",func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		file,err := ctx.FormFile("file")
		if err!=nil{
			ginconfig.Logger_caller("Read picture failed!",err)
			ctx.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}
		app := ctx.PostForm("app")
		// 判断是否允许的图片格式并更改文件名
		dst,err := gincontroller.Change_file_name(file.Filename,name,app) 
		if err!=nil{
			ginconfig.Logger_caller("Add picture failed!",err)
			ctx.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}
		// 保存图片并返回结果
		err = ctx.SaveUploadedFile(file, dst)
		if err!=nil{
			ginconfig.Logger_caller("Add picture failed!",err)
			ctx.JSON(http.StatusInternalServerError,gin.H{"Error":err.Error()})
			return
		}
		ctx.JSON(http.StatusOK,gin.H{"result":"success"})
	})
}