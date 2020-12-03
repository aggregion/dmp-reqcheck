package main

import (
	"context"
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
		Data interface{}
	}

	// JobQueue .
	JobQueue interface {
		// Enqueue puts a job to the queue
		Enqueue(ctx context.Context, job *Job) error

		// Listen retrieves a job from the queue
		Listen(callback func(ctx context.Context, job *Job) error) (UnsubscribeFunc, error)
	}

	localIdempotentHelper struct {
		cache common.Cache
	}

	localJobQueue struct {
		name  string
		queue chan *Job
	}
)

var errQueueBackpressure = errors.New("queue backpressure")

func (lj *localJobQueue) Enqueue(ctx context.Context, job *Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case lj.queue <- job:
		return nil
	}
}

func (lj *localJobQueue) Listen(callback func(ctx context.Context, job *Job) error) (UnsubscribeFunc, error) {
	rootCtx, cancel := context.WithCancel(context.Background())

	waitWorker := sync.WaitGroup{}
	waitWorker.Add(1)
	go func() {
		defer waitWorker.Done()

		var err error
		for job := range lj.queue {
			// log.Debugf("listen and received local job %s", jsonData)
			err = callback(rootCtx, job)
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

// NewLocalJobQueue creates a local job queue
func NewLocalJobQueue(name string, bufferQueueSize int) (JobQueue, error) {
	return &localJobQueue{
		name:  name,
		queue: make(chan *Job, bufferQueueSize),
	}, nil
}
