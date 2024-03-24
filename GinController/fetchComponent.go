package gincontroller

import (
	"encoding/json"
	ginconfig "sing2cat_web/GinConfig"

	"github.com/bitly/go-simplejson"
)

func formatComponent(components []ginconfig.Component) []*simplejson.Json {
	// 将获取到的链接转为json格式返回
	content := make([]*simplejson.Json, len(components))
	for index, component := range components {
		var temp_map map[string]interface{}
		err := json.Unmarshal([]byte(component.Gorm_data),&temp_map)
		if err != nil{
			ginconfig.Logger_caller("Unmarshal failed!",err)
			continue
		}
		// 先变成字典
		component_property := map[string]interface{}{"icon": component.Icon, "url": component.Url, "name": component.Name, "data": temp_map, "class": component.Class}
		// 由字典转为byte
		component_property_byte, _ := json.Marshal(component_property)
		// 最后由simplejson库转为json结构体
		content[index], err = simplejson.NewJson(component_property_byte)
		if err != nil{
			ginconfig.Logger_caller("marshal failed!",err)
			continue
		}
	}
	return content
}
func Fetch_components() ([]*simplejson.Json, error) {
	components := []ginconfig.Component{}
	if err := ginconfig.Db.Select("url", "icon", "name", "gorm_data", "class").Find(&components).Error; err != nil {
		ginconfig.Logger_caller("Get components failed",err)
		return nil, err
	}

	component_json := formatComponent(components)
	return component_json, nil
}