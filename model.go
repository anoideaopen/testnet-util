package utils

type Request struct {
	Args        [][]byte `json:"args"`
	ChaincodeID string   `json:"chaincodeId"`
	Fcn         string   `json:"fcn"`
}

type Response struct {
	BlockNumber      int64  `json:"blockNumber,omitempty"`
	ChaincodeStatus  int64  `json:"chaincodeStatus,omitempty"`
	Payload          []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	TransactionID    string `json:"transactionId,omitempty"`
	TxValidationCode int64  `json:"txValidationCode,omitempty"`
}

type ResponseError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
