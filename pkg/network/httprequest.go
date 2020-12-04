package network

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/aggregion/dmp-reqcheck/pkg/utils"
)

type (
	responseContainer struct {
		Response *http.Response
		Error    error
	}

	httpResponseChannel = chan responseContainer
)

var (
	// ErrHTTPRequestCancel raised when cancel request
	ErrHTTPRequestCancel = errors.New("cancel request")

	// ErrHTTPClientError for client
	ErrHTTPClientError = errors.New("client error")

	// ErrHTTPServerError for server
	ErrHTTPServerError = errors.New("server error")
)

// HTTPRequestAndGetResponse make request with timeout and cancellation.
// Note: Need to manually close response.Body if response != nil
func HTTPRequestAndGetResponse(requestContext context.Context, timeout time.Duration,
	httVerb, url string, body io.Reader, headers map[string][]string, withProxy string) (response *http.Response, err error) {

	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		IdleConnTimeout:       timeout,
		ResponseHeaderTimeout: timeout,
		DisableKeepAlives:     true,
		MaxIdleConns:          1,
	}
	if withProxy == "" {
		tr.Proxy = nil
	}
	if withProxy != ".env" {
		tr.Proxy = http.ProxyURL(utils.MustURLParse(withProxy))
	}

	client := &http.Client{Transport: tr}

	responseChannel := make(httpResponseChannel)
	request, err := http.NewRequest(httVerb, url, body)
	if err != nil {
		return nil, err
	}

	request.Header = headers

	go func() {
		request.Close = true
		response, err := client.Do(request)

		if err == nil {
			if response.StatusCode >= 500 {
				err = ErrHTTPServerError
			} else if response.StatusCode >= 400 {
				err = ErrHTTPClientError
			}
		}

		responseChannel <- responseContainer{
			Response: response,
			Error:    err,
		}
	}()

	ctx, cancel := context.WithTimeout(requestContext, timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		tr.CancelRequest(request)
		tr.CloseIdleConnections()
		<-responseChannel
		close(responseChannel)
		return nil, ctx.Err()
	case r := <-responseChannel:
		close(responseChannel)
		return r.Response, r.Error
	}
}
