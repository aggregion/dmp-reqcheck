package common

import (
	"context"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

// SleepFunction .
type SleepFunction func(ctx context.Context) error

// SleepRandomFunc .
func SleepRandomFunc(waitMin, waitMax time.Duration) SleepFunction {
	rng := waitMax - waitMin
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	interval := time.Duration(float64(waitMin) + float64(rng)*rand.Float64())
	return func(ctx context.Context) error {
		select {
		case <-time.After(interval):
			interval = time.Duration(float64(waitMin) + float64(rng)*rand.Float64())
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// SleepExponentialFunc .
func SleepExponentialFunc(waitInterval time.Duration, multiply float64) SleepFunction {
	interval := waitInterval
	return func(ctx context.Context) error {
		select {
		case <-time.After(interval):
			interval = time.Duration(float64(interval) * multiply)
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// SleepLinearFunc .
func SleepLinearFunc(waitInterval time.Duration) SleepFunction {
	return func(ctx context.Context) error {
		select {
		case <-time.After(waitInterval):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// SleepDummyFunc .
func SleepDummyFunc() SleepFunction {
	return func(ctx context.Context) error {
		return nil
	}
}

// SleepCombineFuncs .
func SleepCombineFuncs(
	firstSleepFunction,
	secondSleepFunction SleepFunction,
	firstAttempts int,
) SleepFunction {
	return func(ctx context.Context) error {
		if firstAttempts > 0 {
			firstAttempts--
			return firstSleepFunction(ctx)
		}
		return secondSleepFunction(ctx)
	}
}

// SleepFuncFactory .
type SleepFuncFactory func() SleepFunction

// RetryMethod .
func RetryMethod(
	ctx context.Context,
	sleep SleepFunction,
	maxAttemptCount int,
	doCallback func(context.Context) error,
) (err error) {
	var sleepErr error
	for {
		err = doCallback(ctx)
		if err == nil {
			return nil
		}
		if ctx.Err() != nil {
			return errors.Wrap(err, ctx.Err().Error())
		}

		if maxAttemptCount == 0 {
			break
		}
		// Infinity retry if pass negative value
		if maxAttemptCount >= 0 {
			maxAttemptCount--
		}

		sleepErr = sleep(ctx)
		if sleepErr != nil {
			return errors.Wrapf(err, "sleep function error: %s", sleepErr.Error())
		}
		if ctx.Err() != nil {
			return errors.Wrap(err, ctx.Err().Error())
		}
	}
	return
}

// RepeatMethod .
func RepeatMethod(
	ctx context.Context,
	sleepFn SleepFunction,
	startImmediately bool,
	method func(ctx context.Context) error,
) (err error) {
	for {
		if !startImmediately {
			err = sleepFn(ctx)
			if err != nil {
				return
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}

		err = method(ctx)
		if err != nil {
			return
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}

		startImmediately = false
	}
}
