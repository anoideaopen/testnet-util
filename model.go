package utils

// Request struct for request to hlf proxy
type Request struct {
	Args        [][]byte `json:"args"`
	ChaincodeID string   `json:"chaincodeId"`
	Fcn         string   `json:"fcn"`
}

// Response struct for response from hlf proxy
type Response struct {
	BlockNumber      int64  `json:"blockNumber,omitempty"`
	ChaincodeStatus  int64  `json:"chaincodeStatus,omitempty"`
	Payload          []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	TransactionID    string `json:"transactionId,omitempty"`
	TxValidationCode int64  `json:"txValidationCode,omitempty"`
}

// ResponseError struct for response error from hlf proxy
type ResponseError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
