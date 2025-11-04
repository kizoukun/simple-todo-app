package mock

var datas = make(map[string]interface{})

func SetData(key string, value interface{}) {
	datas[key] = value
}

func GetData(key string) interface{} {
	if val, exists := datas[key]; exists {
		return val
	}
	return nil
}

func ClearData() {
	datas = make(map[string]interface{})
}
