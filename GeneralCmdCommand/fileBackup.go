package generalcmdcommand

import (
	"fmt"
	"os"
	"path/filepath"
	sing2catconfig "sing2cat_web/Sing2catConfig"
)

func Backup_file(src string,app string) error{
	project_dir,err := sing2catconfig.Get_value("project_dir")
	if err != nil{
		sing2catconfig.Logger_caller("get project_dir failed!",err)
		return err
	}
	_,file_name := filepath.Split(src)
	if _,err := os.Stat(fmt.Sprintf("%s/temp/%s", project_dir,app));err!=nil{
		if os.IsNotExist(err) {
			os.MkdirAll(fmt.Sprintf("%s/temp/%s", project_dir,app),0666)
		}
	}
	// 复制原本的配置文件
	origin_content, err := os.ReadFile(src)
	if err != nil {
		sing2catconfig.Logger_caller("Copy file failed!",err)
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/temp/%s/%s.bak", project_dir,app,file_name), origin_content, 0644)
	if err != nil {
		sing2catconfig.Logger_caller("Copy file failed!",err)
		return err
	}
	// 删除原本的配置文件
	del_err := os.Remove(src)
	if os.IsNotExist(del_err) {
		new_content, err := os.ReadFile(fmt.Sprintf("%s/temp/%s/%s", project_dir,app,file_name))
		if err != nil {
			sing2catconfig.Logger_caller("Read file failed!",err)
			return err
		}
		os.WriteFile(src, new_content, 0644)
	}
	// 将新生成的配置文件复制到目标文件夹
	new_content, err := os.ReadFile(fmt.Sprintf("%s/temp/%s/%s",project_dir,app,file_name))
	if err != nil {
		sing2catconfig.Logger_caller("Read file failed!",err)
		return err
	}
	err = os.WriteFile(src, new_content, 0644)
	if err != nil{
		sing2catconfig.Logger_caller("Move File failed!",err)
		return err
	}
	return nil
}
func Recover_file(src string,app string) error {
	project_dir,err := sing2catconfig.Get_value("project_dir")
	if err != nil{
		sing2catconfig.Logger_caller("get project_dir failed!",err)
		return err
	}
	_,file_name := filepath.Split(src)
	// 删除原本的配置文件
	del_err := os.Remove(src)
	if os.IsNotExist(del_err) {
		backup_content, _ := os.ReadFile(fmt.Sprintf("%s/temp/%s/%s.bak",project_dir,app,file_name))
		os.WriteFile(src, backup_content, 0644)
	}
	// 将之前的配置文件恢复回去
	backup_content, _ := os.ReadFile(fmt.Sprintf("%s/temp/%s/%s.bak", project_dir,app,file_name))
	os.WriteFile(src, backup_content, 0644)
	return nil
}