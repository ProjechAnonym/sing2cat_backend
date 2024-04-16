package main

import (
	_ "embed"
	"fmt"
	"net/http"
	generalcmdcommand "sing2cat_web/GeneralCmdCommand"
	ginconfig "sing2cat_web/GinConfig"
	ginrouter "sing2cat_web/GinRouter"
	middlewarefunc "sing2cat_web/MiddlewareFunc"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func init()  {
	ginconfig.Set_value(ginconfig.Get_gin_dir(),"project_dir")
	ginconfig.Get_core()
	ginconfig.Gin_init()
}
//go:embed static/index.html
var html string
func main() {
	var cronlock sync.RWMutex
	cron := cron.New()
	project_dir,_ := ginconfig.Get_value("project_dir")
	frontend_routes,_ := ginconfig.Get_value("config","routes") 
	Entry_id,_ := cron.AddFunc("30 * * * *",func() {
		status := cronlock.TryLock()
		if status{
			// 成功上锁执行的逻辑
			// 记得解锁
			defer cronlock.Unlock()
			// 运行sing2cat
			generalcmdcommand.Update_config("opt/singbox/config.json","sing-box.service","sing2cat")
		}else{
			return
		}
	})
	cron.Start()
	r := gin.Default()
	r.Use(middlewarefunc.Gin_logger(),middlewarefunc.Gin_recovery(true),cors.New(middlewarefunc.Cors()))
	r.StaticFS("/static",http.Dir(fmt.Sprintf("%s/build/static",project_dir)))
	r.StaticFS("/build",http.Dir(fmt.Sprintf("%s/build",project_dir)))
	for _, route := range frontend_routes.([]interface{}) {
		r.GET(route.(string),func(ctx *gin.Context) {
			ctx.File(fmt.Sprintf("%s/build/index.html",project_dir))
		})
	}
	api := r.Group("/api")
	api.StaticFS("/static",http.Dir(fmt.Sprintf("%s/static",project_dir)))
	ginrouter.Sing2cat_func(api,&Entry_id,cron,&cronlock,"/opt/singbox/config.json","sing-box.service","sing2cat")
	ginrouter.Authentication_router(api,html)
	ginrouter.Add_item(api)
	ginrouter.Fetch_router(api)
	ginrouter.Remove_router(api)
	port,_ := ginconfig.Get_value("config","api","listen")
	r.Run(fmt.Sprintf(":%d",port.(int)))
}