package common

// Request 封装rpc底层封装的请求格式
type Request struct {
	MethodName string
	Args       interface{}
}

// NewRequest 实例化一个新的request结构体指针实例
func NewRequest(methodName string, args interface{}) *Request {
	return &Request{
		MethodName: methodName,
		Args:       args}
}

