package utils

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"strings"
)

func ParseError(err error) ChatErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), "no rows in result set") {
			return ErrorKind(NotFoundError, "no record matching given id")
		}
		return ErrorKind(InternalServerError, fmt.Sprintf("error_utils when trying to save chat: %s", err.Error()))
	}

	switch sqlErr.Number {
	case 1062:
		return ErrorKind(InternalServerError, "title already taken")
	}
	return ErrorKind(InternalServerError, fmt.Sprintf("error_utils when processing request: %s", err.Error()))
}
