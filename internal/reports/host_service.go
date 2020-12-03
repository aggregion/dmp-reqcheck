package reports

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aggregion/dmp-reqcheck/internal/serve"
	"github.com/aggregion/dmp-reqcheck/pkg/common"
	"github.com/aggregion/dmp-reqcheck/pkg/network"
	"github.com/aggregion/dmp-reqcheck/pkg/utils"
)

const (
	// HostAccessibleIntAttr .
	HostAccessibleIntAttr = "accessible"
)

type (
	// HostServiceReport .
	HostServiceReport struct {
		Target     string
		Timeout    time.Duration
		WithProxy  bool
		IsThisHost bool

		accessible int64 `attr:"accessible"`

		httpServer serve.HTTPStubServer
	}
)

var _ = (IReport)((*HostServiceReport)(nil))

// serveHTTP .
func (dr *HostServiceReport) serveHTTP(ctx context.Context, hostURL *url.URL) error {
	path := hostURL.Path
	if path == "" {
		path = "/"
	}
	status, _ := strconv.Atoi(hostURL.Query().Get("response_status"))
	if status == 0 {
		status = 200
	}

	dr.httpServer.Host = fmt.Sprintf("0.0.0.0:%s", hostURL.Port())
	dr.httpServer.ServePath = path
	dr.httpServer.ResponseBody = hostURL.Query().Get("response_body")
	dr.httpServer.ResponseStatus = status

	return dr.httpServer.Start(ctx)
}

// Start .
func (dr *HostServiceReport) startLinux(ctx context.Context) error {
	hostURL := utils.MustURLParse(dr.Target)

	switch hostURL.Scheme {
	case "http":
		return dr.serveHTTP(ctx, hostURL)
	default:
	}

	return errors.New("supported serve http only yet")
}

// Stop .
func (dr *HostServiceReport) stopLinux(ctx context.Context) error {
	return dr.httpServer.Stop(ctx)
}

// Start .
func (dr *HostServiceReport) Start(ctx context.Context) error {
	if !dr.IsThisHost {
		return nil
	}

	return dr.startLinux(ctx)
}

// Stop .
func (dr *HostServiceReport) Stop(ctx context.Context) error {
	if !dr.IsThisHost {
		return nil
	}

	return dr.stopLinux(ctx)
}

func (dr *HostServiceReport) gatherLinux(ctx context.Context) []error {
	dr.accessible = 0

	url := utils.MustURLParse(dr.Target)

	if url.Scheme != "http" {
		return []error{errors.New("supported http only request")}
	}

	timeout := dr.Timeout
	if timeout == 0 {
		timeout = time.Second * 2
	}

	requestMethod := url.Query().Get("method")
	if requestMethod == "" {
		requestMethod = "GET"
	}
	sendBody := url.Query().Get("send_body")
	matchStatusStr := url.Query().Get("match_status")
	if matchStatusStr == "" {
		matchStatusStr = url.Query().Get("response_status")
	}
	if matchStatusStr == "" {
		matchStatusStr = "200"
	}
	matchBodyStr := url.Query().Get("match_body")
	url.Query().Del("match_body")
	if len(matchBodyStr) == 0 {
		matchBodyStr = url.Query().Get("response_body")
	}
	if len(matchBodyStr) == 0 {
		matchBodyStr = ".*"
	}

	var sendBodyStream io.Reader

	if requestMethod != "GET" {
		sendBodyStream = strings.NewReader(sendBody)
	}

	query := url.Query()
	query.Del("method")
	query.Del("send_body")
	query.Del("match_body")
	query.Del("match_status")
	query.Del("response_body")
	query.Del("response_status")
	url.RawQuery = query.Encode()

	common.RetryMethod(ctx, common.SleepExponentialFunc(time.Second, 1.2), 2, func(ctx context.Context) error {
		response, err := network.HTTPRequestAndGetResponse(ctx,
			timeout,
			requestMethod,
			url.String(), sendBodyStream, nil, false)
		if response != nil {
			defer response.Body.Close()
		}

		if err == nil {
			if response == nil {
				return errors.New("response is nil")
			}
			bodyBytes, _ := ioutil.ReadAll(response.Body)
			body := string(bodyBytes)

			if regexp.MustCompile(matchStatusStr).MatchString(fmt.Sprintf("%d", response.StatusCode)) &&
				regexp.MustCompile(matchBodyStr).MatchString(body) {
				dr.accessible = 1
			}
		} else {
			return err
		}

		return nil
	})

	return nil
}

// Gather .
func (dr *HostServiceReport) Gather(ctx context.Context) []error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *HostServiceReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *HostServiceReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
