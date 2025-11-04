package mock

var datas = make(map[string]interface{})

type MockRedis struct {
}

func NewMockRedis() *MockRedis {
	return &MockRedis{}
}

func (m *MockRedis) SetData(key string, value interface{}) {
	datas[key] = value
}

func (m *MockRedis) GetData(key string) interface{} {
	if val, exists := datas[key]; exists {
		return val
	}
	return nil
}

func (m *MockRedis) DeleteData(key string) {
	delete(datas, key)
}
