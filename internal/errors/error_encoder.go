package errors

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
)

func NewHTTPError(code int, field string, detail string) *HTTPError {
	return &HTTPError{
		Code: code,
		Errors: map[string][]string{
			field: {detail},
		},
	}
}

type HTTPError struct {
	Errors map[string][]string `json:"errors"`
	Code   int                 `json:"-"`
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTPError: %d", e.Code)
}

func FromError(err error) *HTTPError {
	if err == nil {
		return nil
	}
	if se := new(HTTPError); errors.As(err, &se) {
		return se
	}
	// 将 kratos 标准错误转换为我们自定义的错误
	if se := new(errors.Error); errors.As(err, &se) {
		return NewHTTPError(int(se.Code), se.Reason, se.Message)
	}
	return NewHTTPError(500, "internal", "error")
}
