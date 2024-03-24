package mergesingboxconfig

import (
	sing2catconfig "sing2cat_web/Sing2catConfig"
)

func format_ruleset_source() ([]map[string]interface{}, error) {
	default_rule_set,err := sing2catconfig.Get_value("template","route","rule_set")
	if err != nil{
		sing2catconfig.Logger_caller("Get default_rule_set failed!",err)
		return nil,err
	}
	custom_rule_set,err := sing2catconfig.Get_value("config","rule_set")
	if err != nil{
		sing2catconfig.Logger_caller("Get custom_rule_set failed!",err)
		return nil,err
	}
	
	// 创建列表用于存储对象
	rule_set_sources := make([]map[string]interface{},len(default_rule_set.([]interface{}))+len(custom_rule_set.([]interface{}))) 
	for _, rule := range(custom_rule_set.([]interface{})){
		// 提取出规则集名称以及路径
		tag := rule.(map[string]interface{})["label"].(string)
		path := rule.(map[string]interface{})["value"].(map[string]interface{})["path"].(string)
		// 创建字典存储值
		rule_set_source := make(map[string]interface{})
		switch rule.(map[string]interface{})["value"].(map[string]interface{})["type"] {
			case "local":
				rule_set_source = map[string]interface{}{"type": "local", "tag": tag, "format": "binary", "path": path}
			case "remote":
				rule_set_source = map[string]interface{}{"type": "remote", "tag": tag, "format": "binary", "url": path, "download_detour": "select", "update_interval": "1d"}
		}
		// 添加新的规则集
		default_rule_set = append(default_rule_set.([]interface{}), rule_set_source)
	}
	// 将其转变为json对象
	for i,rule_set := range(default_rule_set.([]interface{})){
		rule_set_sources[i] = rule_set.(map[string]interface{})
	}
	return rule_set_sources,nil
}

func format_ruleset() ([]map[string]interface{}, error) {
	base_rules,err := sing2catconfig.Get_value("template","route","rules","default")
	if err != nil{
		sing2catconfig.Logger_caller("Get base_rules failed!",err)
		return nil,err
	}
	custom_rules,err := sing2catconfig.Get_value("config","rule_set")
	if err != nil{
		sing2catconfig.Logger_caller("Get custom_rules failed!",err)
		return nil,err
	}
	// 获取分流规则
	shunt_rules,err := sing2catconfig.Get_value("template","route","rules","shunt")
	if err != nil{
		sing2catconfig.Logger_caller("Get shunt_rules failed!",err)
		return nil,err
	}
	rules := make([]map[string]interface{}, len(base_rules.([]interface{}))+len(custom_rules.([]interface{}))+len(shunt_rules.([]interface{})))

	// 此处逻辑于上面相同
	for _, rule := range custom_rules.([]interface{}) {
		tag := rule.(map[string]interface{})["label"].(string)
		switch rule.(map[string]interface{})["value"].(map[string]interface{})["china"] {
			case true:
				base_rules = append(base_rules.([]interface{}), map[string]interface{}{"rule_set": tag, "outbound": "direct"})
			case false:
				base_rules = append(base_rules.([]interface{}), map[string]interface{}{"rule_set": tag, "outbound": tag + "-select"})
		}
	}
	base_rules = append(base_rules.([]interface{}), shunt_rules.([]interface{})...)
	for i,rule_set := range(base_rules.([]interface{})){
		rules[i] = rule_set.(map[string]interface{})
	}
	return rules, nil
}

func Merge_route() (map[string]interface{},error) {
	// 获取规则集源
	rule_set_source,err := format_ruleset_source()
	if err != nil {
		return nil,err		
	}
	// 获取路由规则
	rules,err := format_ruleset()
	if err != nil {
		return nil,err		
	}

	route,err := sing2catconfig.Get_value("template","route")
	if err != nil{
		sing2catconfig.Logger_caller("Marshal route failed!",err)
		return nil,err
	}
	route.(map[string]interface{})["rule_set"] = rule_set_source
	route.(map[string]interface{})["rules"] = rules
	return route.(map[string]interface{}),nil
}