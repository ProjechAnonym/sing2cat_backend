package gincontroller

import (
	"fmt"
	"os"
	ginconfig "sing2cat_web/GinConfig"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Url      []string      `json:"url" yaml:"url"`
	Rule_set []map[string]interface{} `json:"rule_set" yaml:"rule_set"`
}

func Generate_sing2cat_config(config Config) error{
	project_dir,err := ginconfig.Get_value("project_dir")
	if err != nil{
		ginconfig.Logger_caller("Get project dir failed!",err)
	}
	config_byte,err := yaml.Marshal(config)
	if err!=nil{
		return err
	}
	if err = os.WriteFile(fmt.Sprintf("%s/config/sing2cat/config.yaml",project_dir),config_byte,0644);err!=nil{
		return err
	}
	return nil
}
func Interval_spec(spec []string) (string,bool){
	if len(spec) == 1{
		return "",true
	}else if len(spec) == 2{
		return fmt.Sprintf("%s %s * * *",spec[0],spec[1]),false
	}else{
		return fmt.Sprintf("%s %s * * %s",spec[0],spec[1],spec[2]),false
	}
}

