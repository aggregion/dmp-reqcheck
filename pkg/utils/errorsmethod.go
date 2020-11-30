package utils

import (
	"context"
	"fmt"
)

// WrapPanic wraps any method for recovering after panic
func WrapPanic(method func(ctx context.Context)) func(ctx context.Context) error {
	return func(ctx context.Context) (err error) {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("%+v", e)
			}
		}()
		method(ctx)
		return err
	}
}
