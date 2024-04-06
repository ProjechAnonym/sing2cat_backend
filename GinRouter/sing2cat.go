package ginrouter

import (
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	generalcmdcommand "sing2cat_web/GeneralCmdCommand"
	ginconfig "sing2cat_web/GinConfig"
	gincontroller "sing2cat_web/GinController"
	middlewarefunc "sing2cat_web/MiddlewareFunc"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)


func Sing2cat_func(r *gin.RouterGroup, id *cron.EntryID, cron *cron.Cron, lock *sync.RWMutex,app_msg ...string) {
	sing2cat := r.Group("sing2cat")
	sing2cat.Use(middlewarefunc.Jwt_auth())
	sing2cat.POST("/config", func(ctx *gin.Context) {
		secret := ctx.GetString("token")	
		if !gincontroller.Valid_auth(secret) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorization"})
			return
		}
		// 设置Config的变量用于承接前端传递的值
		var config gincontroller.Config
		// 将前端的数据变成json对象
		if err := ctx.BindJSON(&config); err != nil {
			ginconfig.Logger_caller("Marshal json failed!",err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 运行sing2cat,更新配置文件
		// 获得锁的状态,是否成功上锁
		for {
			status := lock.TryLock()
			if status {
				break
			}
		}
		// 成功上锁执行的逻辑
		// 记得解锁
		defer lock.Unlock()
		// 将生成的数据写入sing2cat配置文件
		if err := gincontroller.Generate_sing2cat_config(config); err != nil {
			ginconfig.Logger_caller("Marshal json to yaml failed!",err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 运行sing2cat
		if err := generalcmdcommand.Update_config(app_msg[0],app_msg[1],app_msg[2]);err != nil{
			ginconfig.Logger_caller("Generate sing2cat config file failed!",err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 
		}
		ctx.JSON(http.StatusOK, gin.H{"result": "success"})
	})

	sing2cat.POST("/interval", func(ctx *gin.Context) {
		secret := ctx.GetString("token")
		if !gincontroller.Valid_auth(secret) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorization"})
			return
		}
		// 获取前端的更新时间间隔
		spec := ctx.PostFormArray("time")
		time_spec, disable := gincontroller.Interval_spec(spec)
		// 删除原本的定时任务
		cron.Remove(*id)

		// 添加新的定时任务
		if !disable {
			*id, _ = cron.AddFunc(time_spec, func() {
				// 获得锁的状态,是否成功上锁
				status := lock.TryLock()
				if status {
					// 成功上锁执行的逻辑
					// 记得解锁
					defer lock.Unlock()
					if err := generalcmdcommand.Update_config(app_msg[0],app_msg[1],app_msg[2]);err != nil{
						ginconfig.Logger_caller("Generate sing2cat config file failed!",err)
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return 
					}
				}
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"result": "success"})
	})

	sing2cat.POST("/check",func(ctx *gin.Context) {
		service := ctx.PostForm("service")
		cmd_check := exec.Command("systemctl", "status", service)
		// 获取status的输出结果
		output, err := cmd_check.CombinedOutput()
		if err != nil{
			ginconfig.Logger_caller("Process check command failed!",err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 
		}
		// 判断singbox是否在运行
		if strings.Contains(string(output), "active (running)") {
			ctx.JSON(http.StatusOK, gin.H{"result": "success"})
			return 
		} else {
			msg := fmt.Sprintf("%s is dead",service)
			err = errors.New(msg)
			ginconfig.Logger_caller(msg,err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 
		}
	})
	sing2cat.POST("/restart",func(ctx *gin.Context) {
		service := ctx.PostForm("service")
		if err := generalcmdcommand.Cmd_reboot_service(service);err != nil{
			ginconfig.Logger_caller("restart failed!",err)
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		ctx.JSON(http.StatusOK,gin.H{"result":"success"})
	})
	
}