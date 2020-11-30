package utils

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// DieIfError if error is not empty wrap it in message, print and stop process
func DieIfError(err error, message string) {
	if err == nil {
		return
	}

	if len(message) != 0 {
		err = errors.Wrap(err, message)
	}

	Die(err)
}

// Die prints "last word" and stop process
func Die(lastWord interface{}) {
	fmt.Println(lastWord)
	os.Exit(1)
}
