package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"sync"

	"github.com/aggregion/dmp-reqcheck/pkg/common"
	"github.com/sirupsen/logrus"
)

type (
	// UnsubscribeFunc .
	UnsubscribeFunc func() error

	// Job .
	Job struct {
		Type string
		Data interface{}
	}

	// JobData .
	JobData struct {
		Type string
		Data string
	}

	// JobQueue .
	JobQueue interface {
		// Enqueue puts a job to the queue
		Enqueue(ctx context.Context, job Job) error

		// Listen retrieves a job from the queue
		Listen(callback func(ctx context.Context, jsonData []byte) error) (UnsubscribeFunc, error)
	}

	// IdempotentHelper .
	IdempotentHelper interface {
		Fogot(data ...[]byte)
		IsProduced(data ...[]byte) bool
	}

	localIdempotentHelper struct {
		cache common.Cache
	}

	localJobQueue struct {
		name     string
		queue    chan interface{}
		idemProd IdempotentHelper
	}
)

var errQueueBackpressure = errors.New("queue backpressure")

func getHashKeyForData(data ...[]byte) string {
	hasher := md5.New()
	for _, d := range data {
		hasher.Write(d)
	}
	return string(hasher.Sum(nil))
}

func (pm *localIdempotentHelper) Fogot(data ...[]byte) {
	pm.cache.Delete(getHashKeyForData(data...))
}

func (pm *localIdempotentHelper) IsProduced(data ...[]byte) bool {
	_, exists := pm.cache.GetOrSet(getHashKeyForData(data...), func() (interface{}, error) {
		return true, nil
	})

	return exists
}

func (lj *localJobQueue) Enqueue(ctx context.Context, job Job) error {
	var err error
	var jsonData []byte
	jobData := JobData{Type: job.Type}
	jsonData, err = json.Marshal(job.Data)
	if err != nil {
		return err
	}
	jobData.Data = string(jsonData)
	jsonData, err = json.Marshal(jobData)
	if err != nil {
		return err
	}

	if lj.idemProd.IsProduced(jsonData) {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case lj.queue <- jsonData:
		return nil
	}
}

func (lj *localJobQueue) Listen(callback func(ctx context.Context, jsonData []byte) error) (UnsubscribeFunc, error) {
	rootCtx, cancel := context.WithCancel(context.Background())

	waitWorker := sync.WaitGroup{}
	waitWorker.Add(1)
	go func() {
		defer waitWorker.Done()

		var err error
		for job := range lj.queue {
			jsonData := job.([]byte)
			// log.Debugf("listen and received local job %s", jsonData)
			err = callback(rootCtx, jsonData)
			if err != nil {
				logrus.Errorf("[%s] callback execution failed [%s] [%+v]", lj.name, job, err)
			}
			// if err == errQueueBackpressure {
			// 	// enqueue again for later processing
			// 	go func() {
			// 		select {
			// 		case <-rootCtx.Done():
			// 			return
			// 		case lj.queue <- jsonData:
			// 		}
			// 	}()
			// }
		}
	}()

	return func() error {
		cancel()
		close(lj.queue)
		waitWorker.Wait()
		return nil
	}, nil
}

// NewLocalIdempotentHelper .
func NewLocalIdempotentHelper(maxTTL int) IdempotentHelper {
	return &localIdempotentHelper{
		cache: common.NewTTLCache(int64(maxTTL)),
	}
}

// NewLocalJobQueue creates a local job queue
func NewLocalJobQueue(idemProd IdempotentHelper, name string, bufferQueueSize int) (JobQueue, error) {
	return &localJobQueue{
		idemProd: idemProd,
		name:     name,
		queue:    make(chan interface{}, bufferQueueSize),
	}, nil
}
