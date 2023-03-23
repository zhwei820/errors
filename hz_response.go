package errors

import (
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gin-gonic/gin"
)

type JSONResult struct {
	Code    int64       `json:"code"` // common code please see https://gitlab.matrixport.com/loan/document/-/blob/master/error/error_code.md
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseErr response error, if err is not code error type, default return http.StatusInternalServerError
func HzResponseErr(g *app.RequestContext, err error) {
	if !hzResponseByErr(g, err) {
		HzResponse(g, http.StatusInternalServerError, OkBizCode, nil, err.Error())
	}
}

func hzResponseByErr(g *app.RequestContext, err error) bool {
	if err == nil {
		HzResponse(g, http.StatusOK, OkBizCode, nil, "")
		return true
	}
	message := err.Error()
	err = Cause(err)
	if inner, ok := err.(*CodeError); ok {
		errCode := inner.HTTPCode
		// if custom biz code, should be used
		if inner.BizCode != OkBizCode {
			errCode = inner.BizCode
		}
		HzResponse(g, inner.HTTPCode, errCode, nil, message)
		return true
	}
	return false
}

// ResponseOk response ok
func HzResponseOk(g *app.RequestContext, data interface{}, msg ...string) {
	var s = ""
	if len(msg) > 0 {
		s = msg[0]
	}
	HzResponse(g, http.StatusOK, OkBizCode, data, s)
}

// Response response json, if the above api doesn't satisfy your demands, should be used
func HzResponse(g *app.RequestContext, httpCode, errCode uint32, data interface{}, message string) {
	translatedMsg := getTranslateMsgByLang(string(g.GetHeader("LANGUAGE-TYPE")), errCode)
	if translatedMsg != "" {
		message = translatedMsg
	}
	g.JSON(int(httpCode), gin.H{
		"code":    errCode,
		"message": message,
		"data":    data,
	})
}
