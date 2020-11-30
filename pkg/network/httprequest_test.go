package network

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

const testRequestURL = "http://127.0.0.1:12233/"

var server *http.Server

func bringUpServer(statusCode, delay int) {
	var waitListenAndServe = make(chan struct{})
	var mux = http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		w.WriteHeader(statusCode)
		io.WriteString(w, "Hello, world")
	})

	server = &http.Server{
		Addr:    "127.0.0.1:12233",
		Handler: mux,
	}

	go func() {
		go func() { close(waitListenAndServe) }()
		server.ListenAndServe()
	}()
	<-waitListenAndServe
}

func tearDownServer() {
	server.Close()
}

func TestHttpRequestNormal(t *testing.T) {
	bringUpServer(200, 50)
	defer tearDownServer()
	response, err := HTTPRequestAndGetResponse(context.Background(), time.Minute, "GET", testRequestURL, nil, nil)
	if response == nil {
		t.Fatalf("Expected response not nil, got nil")
	}
	if err != nil {
		t.Fatalf("Expected no errors, got %v", err)
	}

	defer response.Body.Close()

	_, err = ioutil.ReadAll(response.Body)

	if err != nil {
		t.Fatalf("Expected no errors when reading, got %v", err)
	}
}

func TestHttpRequestTimeout(t *testing.T) {
	bringUpServer(200, 50)
	defer tearDownServer()
	response, err := HTTPRequestAndGetResponse(context.Background(), time.Millisecond, "GET", testRequestURL, nil, nil)
	if response != nil {
		t.Fatalf("Expected response nil when timeout, got %v", response)
	}

	if err == nil || !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("Expected DeadlineExceeded, got %v", err)
	}
}

func TestHttpRequestWrongUrl(t *testing.T) {
	bringUpServer(200, 50)
	defer tearDownServer()
	response, err := HTTPRequestAndGetResponse(context.Background(), time.Minute, "GET", "http://127.0.0.1:122337/invalid", nil, nil)
	if response != nil {
		t.Fatalf("Expected response nil for wrong url, got %v", response)
	}

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestHttpRequestCancellation(t *testing.T) {
	bringUpServer(200, 50)
	defer tearDownServer()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	response, err := HTTPRequestAndGetResponse(ctx, time.Minute, "GET", testRequestURL, nil, nil)

	if response != nil {
		t.Fatalf("Expected response nil for wrong url, got %v", response)
	}

	if err == nil || !strings.Contains(err.Error(), "context canceled") {
		t.Fatalf("Expected DeadlineExceeded, got %v", err)
	}
}
