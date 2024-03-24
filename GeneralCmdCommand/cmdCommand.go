package generalcmdcommand

import (
	"errors"
	"fmt"
	"os/exec"
	mergesingboxconfig "sing2cat_web/MergeSingboxConfig"
	sing2catconfig "sing2cat_web/Sing2catConfig"
	"strings"
	"time"
)


func Cmd_reboot_service(service string) error {
	// 停止singbox
	cmd_stop := exec.Command("systemctl", "stop", service)
	if err := cmd_stop.Run(); err != nil {
		return err
	}
	// 隔一秒后启动singbox
	time.Sleep(1 * time.Second)
	cmd_start := exec.Command("systemctl", "start", service)
	if err := cmd_start.Run(); err != nil {
		return err
	}
	// 隔一秒后查看singbox状态
	time.Sleep(1 * time.Second)
	cmd_check := exec.Command("systemctl", "status", service)
	// 获取status的输出结果
	output, err := cmd_check.CombinedOutput()
	if err != nil{
		return err
	}
	// 判断singbox是否在运行
	if strings.Contains(string(output), "active (running)") {
		return nil
	} else {
		msg := fmt.Sprintf("restart %s failed",service)
		err = errors.New(msg)
		return err
	}
}
func Update_config(src string,service string,app string) error{
	if err := mergesingboxconfig.Sing2cat_task();err != nil{
		return err
	}
	err := Backup_file(src,app)
	if err != nil{
		return err
	}
	if err := Cmd_reboot_service(service);err != nil{
		sing2catconfig.Logger_caller("reboot service failed!",err)
		Recover_file(src,app)
		Cmd_reboot_service(service)
		return err
	}
	return nil
}