package reports

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/aggregion/dmp-reqcheck/pkg/common"
	"github.com/aggregion/dmp-reqcheck/pkg/network"
)

const (
	// HTTPStatusIntAttr .
	HTTPStatusIntAttr = "status"
	// HTTPHeaderStrAttrPrefix .
	HTTPHeaderStrAttrPrefix = "header_"
	// HTTPBodyStrAttr .
	HTTPBodyStrAttr = "body"
	// HTTPRequestTimeIntAttr .
	HTTPRequestTimeIntAttr = "req_time"
)

type (
	// HTTPReport .
	HTTPReport struct {
		URL       string
		Method    string
		Headers   map[string][]string
		Body      string
		Timeout   time.Duration
		WithProxy bool

		status  int64                  `attr:"status"`
		headers map[string]interface{} `attrMap:"header_"`
		body    string                 `attr:"body"`
		reqTime int64                  `attr:"req_time"`
	}
)

var _ = (IReport)((*HTTPReport)(nil))

// Start .
func (dr *HTTPReport) Start(ctx context.Context) error {
	return nil
}

// Stop .
func (dr *HTTPReport) Stop(ctx context.Context) error {
	return nil
}

func (dr *HTTPReport) gatherLinux(ctx context.Context) []error {
	dr.status = 0
	dr.headers = make(map[string]interface{})
	dr.body = ""

	if dr.Method == "" {
		dr.Method = "GET"
	}

	var body io.Reader
	if dr.Method != "GET" {
		body = strings.NewReader(dr.Body)
	}

	timeout := dr.Timeout
	if timeout == 0 {
		timeout = time.Second * 2
	}

	common.RetryMethod(ctx, common.SleepExponentialFunc(time.Second, 1.2), 2, func(ctx context.Context) error {
		start := time.Now().UnixNano()

		response, err := network.HTTPRequestAndGetResponse(ctx, timeout, dr.Method, dr.URL, body, dr.Headers, true)
		if response != nil {
			defer response.Body.Close()
		}
		if err == nil {
			if response == nil {
				return errors.New("response is nil")
			}
			dr.status = int64(response.StatusCode)
			for name, values := range response.Header {
				dr.headers[name] = strings.Join(values, " ")
			}

			bodyBytes, _ := ioutil.ReadAll(response.Body)
			dr.body = string(bodyBytes)
		} else {
			return err
		}

		dr.reqTime = (time.Now().UnixNano() - start) / 1000000

		return nil
	})

	return nil
}

// Gather .
func (dr *HTTPReport) Gather(ctx context.Context) []error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *HTTPReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *HTTPReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
