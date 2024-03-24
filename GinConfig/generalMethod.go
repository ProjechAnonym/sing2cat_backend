package ginconfig

import "errors"

func Get_map_value(content map[string]interface{}, keys ...string) (any, error) {
	var value any
	// 逐级获得字典的值
	for i, key := range keys {
		if i != len(keys)-1 {
			content = content[key].(map[string]interface{})
		} else {
			value = content[key]
		}
	}
	if value == nil {
		err := errors.New("fetch map value failed")
		return nil, err
	}
	return value, nil
}