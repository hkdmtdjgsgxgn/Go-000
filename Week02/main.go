package main

import (
	"database/sql"
	"errors"
	"fmt"
)

func main() {
	if err := biz(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println(errors.Unwrap(err).Error())
		}
	}
}
