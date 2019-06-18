package common


type ClientTransfer interface {
	ClientWriteDataByProtocol(data []byte) interface{}
	ClientReadDataByProtocol ()  (uint64, []byte)
}

type ServerTransfer interface {
	ServerReadDataByProtocol () (uint64, []byte, error)
	ServerWriteDataByProtocol(requestID uint64, data []byte)
}