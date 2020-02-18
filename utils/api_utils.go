package utils

import (
	"colleague/kit/errs"
	"fmt"
	"net/http"
	"reflect"

	"github.com/labstack/echo"
)

var ProjectName string = "[Colleague]"

type ApiResult struct {
	Result  interface{} `json:"result"`
	Success bool        `json:"success"`
	Error   ApiError    `json:"error"`
}

type ApiError struct {
	Code    int         `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ArrayResult struct {
	Items      interface{} `json:"items"`
	TotalCount int64       `json:"totalCount"`
}

var (
	// System Error
	ApiErrorSystem             = ApiError{Code: 10001, Message: "System Error"}
	ApiErrorServiceUnavailable = ApiError{Code: 10002, Message: "Service unavailable"}
	ApiErrorRemoteService      = ApiError{Code: 10003, Message: "Remote service error"}
	ApiErrorIPLimit            = ApiError{Code: 10004, Message: "IP limit"}
	ApiErrorPermissionDenied   = ApiError{Code: 10005, Message: "Permission denied"}
	ApiErrorIllegalRequest     = ApiError{Code: 10006, Message: "Illegal request"}
	ApiErrorHTTPMethod         = ApiError{Code: 10007, Message: "HTTP method is not suported for this request"}
	ApiErrorParameter          = ApiError{Code: 10008, Message: "Parameter error"}
	ApiErrorMissParameter      = ApiError{Code: 10009, Message: "Miss required parameter"}
	ApiErrorDB                 = ApiError{Code: 10010, Message: "DB error, please contact the administator"}
	ApiErrorTokenInvaild       = ApiError{Code: 10011, Message: "Token invaild"}
	ApiErrorMissToken          = ApiError{Code: 10012, Message: "Miss token"}
	ApiErrorVersion            = ApiError{Code: 10013, Message: "API version %s invalid"}
	ApiErrorNotFound           = ApiError{Code: 10014, Message: "Resource not found"}
	// Business Error
	ApiErrorUserNotExists = ApiError{Code: 20001, Message: "User does not exists"}
	ApiErrorPassword      = ApiError{Code: 20002, Message: "Password error"}
	ApiErrorWechatContext = ApiError{Code: 20006, Message: "Please operate in the qiye wechat"}

	ApiErrorLogin   = ApiError{Code: 30002, Message: "Login failure"}
	ApiErrorDetails = ApiError{Code: 30003, Message: "Details are as follows"}

	ApiErrorFailedVerify = ApiError{Code: 40005, Message: "failed to verify"}
)

func ReturnApiSucc(ctx echo.Context, status int, totalCount int64, items interface{}) error {
	return ctx.JSON(status, ApiResult{
		Success: true,
		Result:  ArrayResult{TotalCount: totalCount, Items: items},
	})
}
func ReturnResultApiSucc(ctx echo.Context, status int, result interface{}) error {
	return ctx.JSON(status, ApiResult{
		Success: true,
		Result:  result,
	})
}

func ReturnApiWarn(ctx echo.Context, status int, apiError ApiError, err error) error {
	str := ""
	if err != nil {
		str = fmt.Sprint(err)
	}

	return ctx.JSON(status, ApiResult{
		Success: false,
		Error: ApiError{
			Code:    apiError.Code,
			Message: fmt.Sprintf(apiError.Message),
			Details: ProjectName + str,
		},
	})
}

func ReturnApiParameterWarn(c echo.Context, parameters []string) error {
	return c.JSON(http.StatusBadRequest, ApiResult{
		Success: false,
		Error: ApiError{
			Code:    ApiErrorParameter.Code,
			Message: fmt.Sprintf(ApiErrorParameter.Message),
			Details: ProjectName + fmt.Sprint(parameters),
		},
	})
}

func ReturnApiFail(ctx echo.Context, apiError ApiError, err error, v ...interface{}) error {
	status := http.StatusInternalServerError //默认是500错误
	var msg string
	if reflect.TypeOf(err).String() == "*errs.Err" {
		if errResult, ok := err.(*errs.Err); ok {
			status = errResult.Status()
			msg = fmt.Sprint(errResult)
		}
	} else if reflect.TypeOf(err).String() == "*echo.HTTPError" {
		if errResult, ok := err.(*echo.HTTPError); ok {
			status = errResult.Code
			msg = fmt.Sprint(errResult)
		}
	} else {
		msg = fmt.Sprint(err)
	}

	return ctx.JSON(status, ApiResult{
		Success: false,
		Error: ApiError{
			Code:    apiError.Code,
			Message: fmt.Sprintf(apiError.Message, v...),
			Details: ProjectName + msg,
		},
	})
}
