package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	berr "go-scaffold/internal/errors"
	perr "go-scaffold/pkg/errors"
)

// ErrorHandler is HTTP error handler. It sends a JSON response
func ErrorHandler(debug bool, logger *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, ctx echo.Context) {
		if ctx.Response().Committed {
			return
		}

		logger.Error("handle request error", slog.Any("error", err))

		var (
			httpErr *echo.HTTPError
			bErr    *berr.Error

			bc         int
			hintMsg    string
			statusCode int
		)

		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
			bc = berr.ErrInternalError.Code()
			if c, ok := lo.Invert(errHttpStatusCode)[statusCode]; ok {
				bc = c
			}
			hintMsg = fmt.Sprintf("%v", httpErr.Message)
			if une := httpErr.Unwrap(); une != nil {
				err = une
				var be *berr.Error
				if errors.As(une, &be) {
					if be.Unwrap() != nil {
						err = be.Unwrap()
					}
				}
			}
		} else if errors.As(err, &bErr) {
			bc = bErr.Code()
			hintMsg = bErr.Msg()
			statusCode = httpStatusCode(bc)
			if bErr.Unwrap() != nil {
				err = bErr.Unwrap()
			}
		} else {
			de := berr.ErrInternalError
			bc = de.Code()
			hintMsg = hintMessage(de.Label())
			statusCode = httpStatusCode(bc)
		}

		responseBody := NewDefaultBody().
			WithErrNo(bc).
			WithErrMsg(hintMsg)

		if debug {
			errMsg := err.Error()
			if hintMsg != "" {
				errMsg = fmt.Sprintf("%s: %s", hintMsg, err)
			}
			responseBody.WithErrMsg(errMsg)

			stack := perr.ErrorStackTrace(err)
			if stack != nil {
				responseBody.WithStack(stack)
			}
		}

		if ctx.Request().Method == http.MethodHead { // Issue #608
			err = ctx.NoContent(statusCode)
		} else {
			err = ctx.JSON(statusCode, responseBody)
		}
		if err != nil {
			logger.Error("send error response error", slog.Any("error", err))
		}
	}
}

var errHttpStatusCode = map[int]int{
	berr.ErrInternalError.Code():      http.StatusInternalServerError,
	berr.ErrBadCall.Code():            http.StatusBadRequest,
	berr.ErrValidateError.Code():      http.StatusBadRequest,
	berr.ErrInvalidAuthorized.Code():  http.StatusUnauthorized,
	berr.ErrAccessDenied.Code():       http.StatusForbidden,
	berr.ErrResourceNotFound.Code():   http.StatusNotFound,
	berr.ErrCallsTooFrequently.Code(): http.StatusTooManyRequests,
}

func httpStatusCode(c int) int {
	return errHttpStatusCode[c]
}

var errHintMsg = map[string]string{
	berr.ErrInternalError.Label():      "服务器出错",
	berr.ErrBadCall.Label():            "客户端请求错误",
	berr.ErrValidateError.Label():      "参数校验错误",
	berr.ErrInvalidAuthorized.Label():  "未经授权",
	berr.ErrAccessDenied.Label():       "暂无权限",
	berr.ErrResourceNotFound.Label():   "资源不存在",
	berr.ErrCallsTooFrequently.Label(): "请求太频繁",
}

func hintMessage(l string) string {
	return errHintMsg[l]
}
