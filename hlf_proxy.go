package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RequestPrx struct for request to hlf proxy
type RequestPrx struct {
	Args        [][]byte `json:"args"`
	ChaincodeID string   `json:"chaincodeId"`
	Fcn         string   `json:"fcn"`
	Opts        *Options `json:"options,omitempty"`
}

// Options that can be used to specify target endpoints for the transaction invocation.
type Options struct {
	TargetEndpoints []string `json:"targetEndpoints"`
}

// ResponsePrx struct for response from hlf proxy
type ResponsePrx struct {
	BlockNumber      int64  `json:"blockNumber,omitempty"`
	ChaincodeStatus  int64  `json:"chaincodeStatus,omitempty"`
	Payload          []byte `json:"payload,omitempty"`
	TransactionID    string `json:"transactionId,omitempty"`
	TxValidationCode int64  `json:"txValidationCode,omitempty"`
}

// ResponseErrorPrx struct for response error from hlf proxy
type ResponseErrorPrx struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// Invoke ...
func Invoke(ctx context.Context, url, token, cc, fcn string, endpoints []string, args ...string) (*ResponsePrx, error) {
	newCtx, cancel := context.WithTimeout(ctx, InvokeTimeout)
	defer cancel()
	return doRequest(newCtx, url, token, "invoke", cc, fcn, endpoints, args...)
}

// Query ...
func Query(ctx context.Context, url, token, cc, fcn string, endpoints []string, args ...string) (*ResponsePrx, error) {
	newCtx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	return doRequest(newCtx, url, token, "query", cc, fcn, endpoints, args...)
}

//nolint:funlen
func doRequest(
	ctx context.Context,
	url string,
	token, reqType, cc, fcn string,
	endpoints []string,
	args ...string,
) (*ResponsePrx, error) {
	requestData := RequestPrx{
		Args:        AsBytes(args...),
		ChaincodeID: cc,
		Fcn:         fcn,
	}

	if len(endpoints) != 0 {
		requestData = RequestPrx{
			Args:        AsBytes(args...),
			ChaincodeID: cc,
			Fcn:         fcn,
			Opts: &Options{
				TargetEndpoints: endpoints,
			},
		}
	}

	reqBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", url, reqType), bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("http new request: %w", err)
	}

	req.Header.Add("authorization", fmt.Sprintf("Basic %s", token))
	req.Header.Add("content-type", "application/json")

	httpResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client do: %w", err)
	}

	defer func() {
		clErr := httpResponse.Body.Close()
		if clErr != nil {
			fmt.Printf("body close error: %v\n", err)
		}
	}()
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		responseError := &ResponseErrorPrx{}
		err = json.Unmarshal(body, responseError)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal: %w", err)
		}

		return nil, fmt.Errorf(responseError.Message)
	}

	var resp ResponsePrx
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return &resp, nil
}
