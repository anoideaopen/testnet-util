package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// HTTPClient struct
// url - domain and port for observer service,
// example http://localhost:3335/api without '/' on the end the string.
type HTTPClient struct {
	url string
}

func NewHTTPClient(url ...string) *HTTPClient {
	if len(url) == 1 {
		return &HTTPClient{
			url: url[0],
		}
	}
	return &HTTPClient{
		url: GetEnv(ObserverAPIURL, "http://localhost:3305"),
	}
}

func (o *HTTPClient) Post(t provider.T, apiPath string, v any) ([]byte, int) {
	var (
		resp         *http.Response
		body         []byte
		sc           int
		client       http.Client
		u            string
		bytesRequest []byte
		err          error
	)

	t.WithNewStep("Prepare http client", func(sCtx provider.StepCtx) {
		bytesRequest, err = json.Marshal(v)
		t.Require().NoError(err)
		t.WithNewAttachment("POST data", allure.Text, bytesRequest)
		requestTimeout, err := strconv.Atoi(GetEnv("REQUEST_TIMEOUT", "1"))
		t.Assert().NoError(err)
		client = http.Client{
			Timeout: time.Duration(requestTimeout) * time.Second,
		}
		u = o.PrepareURL(t, apiPath)
	})

	t.WithNewStep("POST to: "+u, func(sCtx provider.StepCtx) {
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			u,
			bytes.NewReader(bytesRequest),
		)
		resp, err = client.Do(req)
		t.Require().NoError(err)
		body, err = io.ReadAll(resp.Body)
		t.Require().NoError(err)
		sc = resp.StatusCode
		defer func() {
			err = resp.Body.Close()
			t.Assert().NoError(err)
		}()
	})
	return body, sc
}

func (o *HTTPClient) Get(t provider.T, apiPath string) ([]byte, int) {
	var (
		body   []byte
		sc     int
		u      string
		client http.Client
	)

	t.WithNewStep("Prepare http client", func(sCtx provider.StepCtx) {
		requestTimeout, err := strconv.Atoi(GetEnv("REQUEST_TIMEOUT", "1"))
		sCtx.Assert().NoError(err)
		client = http.Client{
			Timeout: time.Duration(requestTimeout) * time.Second,
		}
		u = o.PrepareURL(t, apiPath)
	})

	t.WithNewStep("GET to: "+u, func(sCtx provider.StepCtx) {
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			u,
			bytes.NewReader(body),
		)
		resp, err := client.Do(req)
		t.Require().NoError(err)
		body, err = io.ReadAll(resp.Body)
		t.Require().NoError(err)
		sc = resp.StatusCode
		defer func() {
			err = resp.Body.Close()
			t.Assert().NoError(err)
		}()
	})
	return body, sc
}

func (o *HTTPClient) PrepareURL(t provider.T, apiPath string) string {
	u, err := url.Parse(o.url)
	t.Require().NoError(err)
	u.Path = path.Join(u.Path, apiPath)
	return u.String()
}
