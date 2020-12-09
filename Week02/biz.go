package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/wedojava/Go-000/Week02/dao"
)

func biz() error {
	if err := dao.Dao(); err != nil {
		return errors.Wrap(err, "biz error")
	}
	fmt.Println("biz done.")
	return nil
}
