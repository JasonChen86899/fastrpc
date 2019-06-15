package common

import (
	"encoding/json"
)

// EdCode 编解码接口
type EdCode interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

// JSONEdCode 创建编解码接口的json序列化与反序列化实现
type JSONEdCode int

func (edcode JSONEdCode) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (edcode JSONEdCode) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// GetEdcode 默认采用json编码
func GetEdcode() (EdCode, error) {
	return *new(JSONEdCode), nil
}
