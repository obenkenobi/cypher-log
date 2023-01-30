package utils

import (
	"errors"
	"github.com/akrennmair/slice"
	"strings"
)

func ConcatErrors(errs ...error) error {
	errs = slice.Filter(errs, func(err error) bool { return err != nil })
	if len(errs) <= 0 {
		return nil
	}
	errStrings := slice.Map(errs, error.Error)
	return errors.New(strings.Join(errStrings, " : "))
}
