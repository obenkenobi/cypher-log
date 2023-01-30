package utils

import "fmt"

func ConcatErrors(err1 error, err2 error) error {
	if err1 == nil && err2 == nil {
		return nil
	} else if err1 == nil {
		return err2
	} else if err2 == nil {
		return err1
	} else {
		return fmt.Errorf("%v:%v", err1.Error(), err2.Error())
	}
}
