package serve

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// HTTPStubServer .
type HTTPStubServer struct {
	Host string

	ServePath      string
	ResponseBody   string
	ResponseStatus int

	waitListenAndServe chan struct{}
}

// ServeHTTP .
func (dr *HTTPStubServer) serveHTTP(ctx context.Context, host string) error {
	var mux = http.NewServeMux()

	mux.HandleFunc(dr.ServePath, func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(dr.ResponseStatus)
		io.WriteString(w, dr.ResponseBody)
	})

	var server *http.Server

	server = &http.Server{
		Addr:    dr.Host,
		Handler: mux,
	}

	dr.waitListenAndServe = make(chan struct{})
	var err error
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			fmt.Printf("busy %s\n skip it\n", server.Addr)
			close(dr.waitListenAndServe)
		}
	}()

	go func() {
		if dr.waitListenAndServe == nil {
			return
		}
		select {
		case <-dr.waitListenAndServe:
		case <-ctx.Done():
			close(dr.waitListenAndServe)
		}

		dr.waitListenAndServe = nil

		server.Close()
	}()

	return nil
}

// Start .
func (dr *HTTPStubServer) Start(ctx context.Context) error {
	return dr.serveHTTP(ctx, dr.Host)
}

// Stop .
func (dr *HTTPStubServer) Stop(ctx context.Context) error {
	if dr.waitListenAndServe != nil {
		close(dr.waitListenAndServe)
	}

	return nil
}
