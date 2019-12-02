package error_formats

import (
	"efficient-api/utils/error_utils"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"strings"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	ErrorNoRows = "no rows in result set"
)


func ParseError(err error) error_utils.MessageErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return error_utils.NewNotFoundError("no record matching given id")
		}
		return error_utils.NewInternalServerError(fmt.Sprintf("error when trying to save message: %s", err.Error()))
	}
	switch sqlErr.Number {
	case 1062:
		return error_utils.NewInternalServerError("title already taken")
	}
	return error_utils.NewInternalServerError(fmt.Sprintf("error when processing request: %s", err.Error()))
}