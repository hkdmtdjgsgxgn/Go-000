package main

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/wedojava/Go-000/Week02/dao"
)

func biz() error {
	if err := dao.Dao(); errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(err, "biz error")
	}
	fmt.Println("biz done.")
	return nil
}
