package gincontroller

import (
	"fmt"
	"os"
	ginconfig "sing2cat_web/GinConfig"
)

func Delete_component(name string,app string) error {
	var component ginconfig.Component
	if err := ginconfig.Db.Select("icon_path").Where("class = ? AND name = ?", app, name).First(&component).Error; err != nil {
		ginconfig.Logger_caller("Connect to database failed!",err)
		return err
	}
	// 判断拼接成的uri和数据库中是否一致
	if component.Icon_path != ""{
		project_dir,err := ginconfig.Get_value("project_dir")
		if err != nil{
			ginconfig.Logger_caller("Get project_dir failed!",err)
			return err
		}
		// 删除图片
		if err := os.Remove(fmt.Sprintf("%s//%s", project_dir, component.Icon_path)); err != nil {
			ginconfig.Logger_caller("Delete pic failed!",err)
			return err
		}
	}
	// 删除该用户指定的网址
	if err := ginconfig.Db.Where("class = ? AND  name = ?", app, name).Delete(&ginconfig.Component{}).Error; err != nil {
		ginconfig.Logger_caller("Connect to database failed!",err)
		return err
	}
	return nil
}