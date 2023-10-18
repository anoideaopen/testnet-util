package utils

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

const (
	// HlfProxyURL - domain and port for hlf proxy service, example http://localhost:9001 without '/' on the end the string
	HlfProxyURL = "HLF_PROXY_URL"
	// HlfProxyAuthToken - support Basic Auth with auth token
	HlfProxyAuthToken = "HLF_PROXY_AUTH_TOKEN" //nolint:gosec
	// FiatIssuerPrivateKey - issuer private key ed25519 in base58 check
	FiatIssuerPrivateKey = "FIAT_ISSUER_PRIVATE_KEY"
	// BatchTransactionTimeout - common time execution of following process
	// robot - defaultBatchLimits.batchTimeoutLimit
	// Time batch execute by robot
	BatchTransactionTimeout = 2 * time.Second
	// ObserverAPIURL - domain and port for observer service, example http://localhost:3335/api without '/' on the end the string
	ObserverAPIURL = "OBSERVER_API_URL"
	// CorrectNodeName Name of any node from stand
	CorrectNodeName = "CORRECT_NODE_NAME"
	// DefaultSwapHash - default swap hash
	DefaultSwapHash = "7d4e3eec80026719639ed4dba68916eb94c7a49a053e05c8f9578fe4e5a3d7ea" // #nosec G101
	// DefaultSwapKey - default swap key
	DefaultSwapKey = "12345"
	// InvokeTimeout sets timeout for invoke method operations
	InvokeTimeout = 10 * time.Second
	// QueryTimeout sets timeout for query method operations
	QueryTimeout = 10 * time.Second
	// QueryTimeout sets timeout for query method operations
	MoreNonceTTL = 11 * time.Second
)

func AsBytes(args ...string) [][]byte {
	bytes := make([][]byte, len(args))
	for i, arg := range args {
		bytes[i] = []byte(arg)
	}
	return bytes
}

func GetNonce() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}

// GetEnv return env if found or defaultValue if not found
func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func CheckStatusCode(t provider.T, expectedResponseCode int, actualResponseCode int) {
	t.WithNewStep("Checking that SC is "+strconv.Itoa(expectedResponseCode), func(sCtx provider.StepCtx) {
		sCtx.Require().Equal(expectedResponseCode, actualResponseCode)
	})
}

func FillStructFromBody(t provider.T, body []byte, tt any) {
	err := json.Unmarshal(body, &tt)
	t.Require().NoError(err)
}

func ConvertPemTox509certificate(bytes []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("PEM decoding resulted in an empty block")
	}
	// Important: This method looks very similar to getCertFromPem(idBytes []byte) (*x509.Certificate, error)
	// But we:
	// 1) Must ensure PEM block is of type CERTIFICATE or is empty
	// 2) Must not replace getCertFromPem with this method otherwise we will introduce
	//    a change in validation logic which will result in a chain fork.
	if block.Type != "CERTIFICATE" && block.Type != "" {
		return nil, errors.New("pem type is " + block.Type + ", should be 'CERTIFICATE' or missing")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
