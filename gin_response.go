package errors

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReturnData struct {
	Code    uint32 `json:"code"`
	Data    interface{}
	Message string `json:"message"`
}

const (
	OkBizCode uint32 = 0
)

// ResponseErr response error, if err is not code error type, default return http.StatusInternalServerError
func ResponseErr(g *gin.Context, err error) {
	if !responseByErr(g, err) {
		Response(g, http.StatusInternalServerError, OkBizCode, nil, err.Error())
	}
}

func responseByErr(g *gin.Context, err error) bool {
	if err == nil {
		Response(g, http.StatusOK, OkBizCode, nil, "")
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
		Response(g, inner.HTTPCode, errCode, nil, message)
		return true
	}
	return false
}

// ResponseOk response ok
func ResponseOk(g *gin.Context, data interface{}, msg ...string) {
	var s = ""
	if len(msg) > 0 {
		s = msg[0]
	}
	Response(g, http.StatusOK, OkBizCode, data, s)
}

// Response response json, if the above api doesn't satisfy your demands, should be used
func Response(g *gin.Context, httpCode, errCode uint32, data interface{}, message string) {
	translatedMsg := getTranslateMsg(g, errCode)
	if translatedMsg != "" {
		message = translatedMsg
	}
	g.JSON(int(httpCode), gin.H{
		"code":    errCode,
		"message": message,
		"data":    data,
	})
}

func getTranslateMsg(g *gin.Context, bizCode uint32) string {
	bizCodeStr := strconv.FormatInt(int64(bizCode), 10)
	translated := TranslateWithConvertLan(g.GetHeader("LANGUAGE-TYPE"), bizCodeStr)
	if translated == bizCodeStr {
		return ""
	}
	return translated
}

func getTranslateMsgByLang(lang string, bizCode uint32) string {
	bizCodeStr := strconv.FormatInt(int64(bizCode), 10)
	translated := TranslateWithConvertLan(lang, bizCodeStr)
	if translated == bizCodeStr {
		return ""
	}
	return translated
}
