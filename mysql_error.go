package errors

import (
	innerErr "errors"
	"net/http"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc/codes"
	gormv2 "gorm.io/gorm"
)

// IsDuplicateError is duplicate error
func IsDuplicateError(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); !ok {
		return false
	} else if mysqlErr.Number == 1062 {
		return true
	}
	return false
}

const (
	MysqlErrorBizCode uint32 = 1200001301
)

// V2MysqlErr gorm version 2 mysql error
func V2MysqlErr(err error) error {
	if err == nil {
		return nil
	}
	if innerErr.Is(err, gormv2.ErrRecordNotFound) {
		return NewNotFound(err, "")
	}
	if IsDuplicateError(err) {
		return NewAlreadyExists(err, "")
	}
	return &CodeError{Err: wrap(err, " mysql error", ""), Code: codes.Internal, HTTPCode: http.StatusInternalServerError, BizCode: MysqlErrorBizCode}
}

func init() {
	strMysqlErrorBizCode := strconv.Itoa(int(MysqlErrorBizCode))
	RegisterTranslate([]TransInfo{
		{
			Tag: EnUs,
			Key: strMysqlErrorBizCode,
			Msg: "network error",
		},
		{
			Tag: ZhCn,
			Key: strMysqlErrorBizCode,
			Msg: "网络错误",
		},
		{
			Tag: ZhTW,
			Key: strMysqlErrorBizCode,
			Msg: "網絡錯誤",
		},
		{
			Tag: RuRu,
			Key: strMysqlErrorBizCode,
			Msg: "network error",
		},
	})
}
