package ginconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/huandu/go-clone"
	"github.com/spf13/viper"
)

var global_map = make(map[string]interface{})

func Get_value(keys ...string) (any, error) {
	var value any
	result := clone.Clone(global_map).(map[string]interface{})
	for i, key := range keys {
		if i != len(keys)-1 {
			temp_result := result[key]
			if temp_result == nil{
				msg := fmt.Sprintf("key %s in Global_vars not found",key)
				err := errors.New(msg)
				return nil,err
			}
			result = temp_result.(map[string]interface{})
		} else {
			value = result[key]
			// 判断result是否为空,为空则报错
		}
	}
	if value != nil {
		return value, nil
	} else {
		err := "key in Global_vars not found"
		return value, errors.New(err)
	}
}
func Set_value(value interface{},keys ...string) error{
	result := global_map
	for i, key := range keys {
		if i != len(keys)-1 {
			temp_result := result[key]
			if temp_result == nil{
				msg := fmt.Sprintf("key %s in Global_vars not found",key)
				err := errors.New(msg)
				return err
			}
			result = temp_result.(map[string]interface{})
		} else {
			result[key] = value
		}
	}
	return nil
}

func Get_gin_dir() string {
	base_dir := filepath.Dir(os.Args[0])
	// base_dir := "E:/Myproject/sing2cat_web"
	return base_dir
}
func get_gin_config(file string) error {
	// 获取项目目录路径,获取失败直接panic退出该进程
	project_dir, err := Get_value("project_dir")
	if err != nil {
		Logger_caller(fmt.Sprintf("Get %s Dir failed!", file), err)
		return err
	}
	// 读取配置文件,读取错误则panic退出该进程
	viper.SetConfigFile(fmt.Sprintf("%s/config/%s.yaml", project_dir, file))
	err = viper.ReadInConfig()
	if err != nil {
		Logger_caller(fmt.Sprintf("Read %s failed!", file), err)
		return err
	}
	Set_value(viper.AllSettings(),file)
	return nil
}

func Gin_init() {
	get_gin_config("config")
	Get_database()
	// 日志记录
	Logger_caller("Initial completed!", nil)
}