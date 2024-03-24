package gincontroller

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	ginconfig "sing2cat_web/GinConfig"
	"time"

	"github.com/tidwall/buntdb"
	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v3"
)



func encrypto(password string) string {
	password_byte := []byte(password)
	md5 := md5.New()
	md5.Write(password_byte)
	return hex.EncodeToString(md5.Sum(nil))
}
func Valid_auth(secret string) bool {
	token, err := ginconfig.Get_value("config", "jwt", "key")
	if err != nil {
		ginconfig.Logger_caller("Get sing-box token failed!", err)
	}
	return secret == token

}

func Send_reset_url(email string,html string) error{
	
	// 获取smtp邮件服务器配置

	smtp_config,err := ginconfig.Get_value("config","smtp")
	if err != nil{
		ginconfig.Logger_caller("Get smtp failed!",err)
		return err
	}
	smtp_host := smtp_config.(map[string]interface{})["host"]
	smtp_port := smtp_config.(map[string]interface{})["port"].(int)
	smtp_username := smtp_config.(map[string]interface{})["username"]
	smtp_password := smtp_config.(map[string]interface{})["password"]
	dial := gomail.NewDialer(
		smtp_host.(string),
		smtp_port,
		smtp_username.(string),
		smtp_password.(string),
	)
	

	// 生成token
	secret_key := fmt.Sprintf("%06v",rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000) )
	token := encrypto(secret_key)
	// 写入BuntDb
	
	if err = ginconfig.BuntDb.Update(func(tx *buntdb.Tx) error {
		_,_,err:= tx.Set(email,token,&buntdb.SetOptions{Expires: true, TTL: 5*time.Minute})
		if err != nil{
			return err
		}
		return nil
	});err != nil{
		ginconfig.Logger_caller("Set reset token to BuntDb failed!",err)
		return err
	}
	// 发送给邮箱的链接
	// 发送邮箱的设置
	msg := gomail.NewMessage()
	msg.SetHeader("From",smtp_username.(string))
	msg.SetHeader("To",email)
	msg.SetHeader("Subject", "找回密码")
	msg.SetBody("text/html", fmt.Sprintf(html,secret_key)) 
	if err := dial.DialAndSend(msg); err != nil {
		ginconfig.Logger_caller("Send email failed!",err)
		return err
	}
	return nil
}
func Edit_password(email string,password string,captcha string) (bool,error){
	// 判断token是否有效
	var valid bool
	// 获取密码byte
	captcha_crypto := encrypto(captcha)
	// 从redis中获取token
	if err := ginconfig.BuntDb.View(func(tx *buntdb.Tx) error {
		value,err := tx.Get(email)
		if err != nil{
			ginconfig.Logger_caller("Get token failed!",err)
			return err
		}
		valid = value == captcha_crypto
		return nil
	});err != nil{
		return false,err
	}
	if valid{
		ginconfig.BuntDb.Update(func(tx *buntdb.Tx) error {
			if _,err := tx.Delete(email); err != nil{
				ginconfig.Logger_caller("Delete item from buntdb failed!",err)
				return err
			}
			return nil
		})
		err := ginconfig.Set_value(password,"config","jwt","key")
		if err != nil{
			ginconfig.Logger_caller("Set config JWT key failed!",err)
			return false,err
		}
		config,err := ginconfig.Get_value("config")
		if err != nil{
			ginconfig.Logger_caller("Get config failed!",err)
			return false,err
		}
		config_byte,err := yaml.Marshal(config)
		if err != nil{
			ginconfig.Logger_caller("Marshal config to yaml failed!",err)
			return false,err
		}
		project_dir,err := ginconfig.Get_value("project_dir")
		if err != nil{
			ginconfig.Logger_caller("Get project dir failed!",err)
			return false,err
		}
		src := fmt.Sprintf("%s/config/config.yaml",project_dir) 
		err = os.Remove(src + "")
		if os.IsNotExist(err) {
			if err = os.WriteFile(src, config_byte, 0644);err != nil{
				ginconfig.Logger_caller("Write File failed!",err)
				return false,err
			}
			return true,nil
		}
		if err = os.WriteFile(src, config_byte, 0644);err != nil{
			ginconfig.Logger_caller("Write File failed!",err)
			return false,err
		}
	}
	return valid,nil
}