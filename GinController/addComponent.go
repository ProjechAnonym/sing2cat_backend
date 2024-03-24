package gincontroller

import (
	"errors"
	"fmt"
	"path/filepath"
	ginconfig "sing2cat_web/GinConfig"
)

func Change_file_name(file_name string, component_name string, app string) (string, error) {
	project_dir,err := ginconfig.Get_value("project_dir")
	if err != nil{
		ginconfig.Logger_caller("Get project_dir failed!",err)
		return "Error",err
	}
	// 允许的文件后缀
	types := []string{".jpg", ".jpeg", ".png", ".svg", ".ico"}
	// 获取文件类型
	file_type := filepath.Ext(file_name)

	// 遍历允许类型查看是否允许类型
	for index, allow_type := range types {

		if file_type == allow_type {
			break
		} else if index == len(types)-1 {
			return "Unaccepted", errors.New("不是允许的格式")
		}
	}
	// 更新新的文件名
	new_file := component_name + file_type
	file_path := fmt.Sprintf("%s/static/%s/%s", project_dir, app,new_file)
	return file_path, nil
}