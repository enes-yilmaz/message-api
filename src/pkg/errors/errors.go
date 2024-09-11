package errors

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

const messageApi string = "message-api"
const validatorOp string = "message-api.handlers.validator"
const repoOp string = "message-api.repository"

var (
	ValidatorError = New(validatorOp, "", 10400, http.StatusBadRequest)
)

// Pagination
var (
	LimitLessThanZeroError  = New(validatorOp, "Limit value must be greater than zero", 4001, http.StatusBadRequest)
	OffsetLessThanZeroError = New(validatorOp, "Offset value must be greater than zero", 4002, http.StatusBadRequest)
	LimitParsingError       = New(validatorOp, "Limit value couldn't parsed", 4002, http.StatusBadRequest)
	OffsetParsingError      = New(validatorOp, "Offset value couldn't parsed", 4002, http.StatusBadRequest)
)

type Error struct {
	Public     PublicError
	StatusCode int
	Internal   error
	Args       interface{}
}

type PublicError struct {
	Op        string
	Desc      string
	ErrorCode int
}

func (e *Error) Error() string {
	return fmt.Sprintf("Operation: %s, Description: %s, ErrorCode: %d, Internal: %v , Args: %v", e.Public.Op, e.Public.Desc, e.Public.ErrorCode, e.Internal, e.Args)
}

func New(op string, desc string, errorCode int, statusCode int) *Error {
	return &Error{Public: PublicError{
		Op:        op,
		Desc:      desc,
		ErrorCode: errorCode,
	}, StatusCode: statusCode}
}

func (e *Error) WrapDesc(desc string) *Error {
	return &Error{Public: PublicError{
		Op:        e.Public.Op,
		Desc:      desc,
		ErrorCode: e.Public.ErrorCode,
	},
		StatusCode: e.StatusCode,
	}
}

func (e *Error) Wrap(err error, args ...interface{}) *Error {
	if err == nil {
		return nil
	}

	return &Error{Public: PublicError{
		Op:        e.Public.Op,
		Desc:      e.Public.Desc,
		ErrorCode: e.Public.ErrorCode,
	},
		StatusCode: e.StatusCode,
		Internal:   err,
		Args:       args,
	}
}

func (e Error) ToResponse(c echo.Context) error {
	return c.JSON(e.StatusCode, e.Public)
}

func (e *Error) Log() {
	fields := logrus.Fields{
		"StatusCode": e.StatusCode,
		"Op":         e.Public.Op,
		"ErrorCode":  e.Public.ErrorCode,
		"Args":       e.Args,
		"Internal":   e.Internal,
	}
	if e.StatusCode >= 400 && e.StatusCode < 500 {
		logrus.WithFields(fields).Info(e.Public.Desc)
	} else if e.StatusCode >= 500 {
		logrus.WithFields(fields).Error(e.Public.Desc)
	}
}
