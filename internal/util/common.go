package util

import (
	"errors"
	"github.com/samber/lo"
)

func MergeErrors(errList []error) error {
	errList = lo.Filter(
		errList, func(item error, _ int) bool {
			return item != nil
		},
	)
	if len(errList) > 0 {
		msg := lo.Reduce(
			errList, func(agg string, err error, i int) string {
				return agg + "\n" + err.Error()
			}, "",
		)
		return errors.New(msg)
	}
	return nil
}
