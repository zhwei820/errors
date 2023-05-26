package errors

import (
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/grpc/codes"
)

type ErrCode uint32

const (
	ErrCodeBadRequest          ErrCode = 12000003
	ErrCodeNotFound            ErrCode = 12000005
	ErrCodeConflict            ErrCode = 12000006
	ErrCodeForbidden           ErrCode = 12000007 // 403
	ErrCodePreconditionFailed  ErrCode = 12000009
	ErrCodeNotImplemented      ErrCode = 12000012
	ErrCodeInternalServerError ErrCode = 12000013
	ErrCodeServiceUnavailable  ErrCode = 12000014
	ErrCodeUnauthorized        ErrCode = 12000016 // 401
)

func (c ErrCode) String() string {
	return strconv.Itoa(int(c))
}
func (c ErrCode) Int() uint32 {
	return uint32(c)
}

func InitI18n() {
	RegisterI18n([]TransInfo{
		{
			Lang: EnUs,
			Key:  ErrCodeBadRequest.String(),
			Msg:  "parameter error",
		},
		{
			Lang: ZhCn,
			Key:  ErrCodeBadRequest.String(),
			Msg:  "参数错误",
		},
		// =================================================================
		{
			Lang: EnUs,
			Key:  ErrCodeNotFound.String(),
			Msg:  "record not found",
		},
		{
			Lang: ZhCn,
			Key:  ErrCodeNotFound.String(),
			Msg:  "未找到",
		},
		// =================================================================
		{
			Lang: EnUs,
			Key:  ErrCodeConflict.String(),
			Msg:  "内部错误",
		},
		{
			Lang: ZhCn,
			Key:  ErrCodeConflict.String(),
			Msg:  "内部错误",
		},
		// =================================================================
		{
			Lang: EnUs,
			Key:  ErrCodeForbidden.String(),
			Msg:  "forbidden",
		},
		{
			Lang: ZhCn,
			Key:  ErrCodeForbidden.String(),
			Msg:  "没有权限,禁止访问",
		},
		// =================================================================
		{
			Lang: EnUs,
			Key:  ErrCodePreconditionFailed.String(),
			Msg:  "precondition error",
		},
		{
			Lang: ZhCn,
			Key:  ErrCodePreconditionFailed.String(),
			Msg:  "前置条件错误",
		},
		// =================================================================
		{
			Lang: EnUs,
			Key:  ErrCodeNotImplemented.String(),
			Msg:  "not implemented",
		},
		{
			Lang: ZhCn,
			Key:  ErrCodeNotImplemented.String(),
			Msg:  "未实现",
		},
		// =================================================================
		{
			Lang: EnUs,
			Key:  ErrCodeInternalServerError.String(),
			Msg:  "internal server error",
		},
		{
			Lang: ZhCn,
			Key:  ErrCodeInternalServerError.String(),
			Msg:  "内部错误,请稍后重试,或者联系管理员",
		},
		// =================================================================
		{
			Lang: EnUs,
			Key:  ErrCodeServiceUnavailable.String(),
			Msg:  "service unavailable",
		},
		{
			Lang: ZhCn,
			Key:  ErrCodeServiceUnavailable.String(),
			Msg:  "服务不可用",
		},
		// =================================================================
		{
			Lang: EnUs,
			Key:  ErrCodeUnauthorized.String(),
			Msg:  "Unauthorized",
		},
		{
			Lang: ZhCn,
			Key:  ErrCodeUnauthorized.String(),
			Msg:  "未登录",
		},
	})
}

// =================================================================

// wrap is a helper to construct an *wrapper.
func wrap(err error, format, suffix string, args ...interface{}) Err {
	newErr := Err{
		message:  fmt.Sprintf(format+suffix, args...),
		previous: err,
	}
	newErr.SetLocation(2)
	return newErr
}

// CodeError represents an error on timeout.
type CodeError struct {
	Err
	Code     codes.Code
	HTTPCode uint32
	BizCode  uint32 // custom biz code
}

// GetGRPCCode return grpc code
func (c *CodeError) GetGRPCCode() codes.Code {
	return c.Code
}

// GetHTTPCode get http code
func (c *CodeError) GetHTTPCode() uint32 {
	return c.HTTPCode
}

// GetBizCode get biz code
func (c *CodeError) GetBizCode() uint32 {
	return c.BizCode
}

// NewCodeError new code error
func NewCodeError(code codes.Code, httpCode, bizCode uint32) error {
	return &CodeError{Err: wrap(nil, "", ""), Code: code, HTTPCode: httpCode, BizCode: bizCode}
}

// NewCodeErrorf new code errorf
func NewCodeErrorf(code codes.Code, httpCode, bizCode uint32, format string, args ...interface{}) error {
	return &CodeError{Err: wrap(nil, format, "", args...), Code: code, HTTPCode: httpCode, BizCode: bizCode}
}

// NotValidf returns an error which satisfies IsNotValid().
func NotValidf(format string, args ...interface{}) error {
	return &CodeError{Err: wrap(nil, format, " "+http.StatusText(http.StatusBadRequest), args...), Code: codes.InvalidArgument, HTTPCode: http.StatusBadRequest, BizCode: ErrCodeBadRequest.Int()}
}

// NewNotValid returns an error which wraps err and satisfies IsNotValid().
func NewNotValid(err error, msg string) error {
	return &CodeError{Err: wrap(err, msg, ""), Code: codes.InvalidArgument, HTTPCode: http.StatusBadRequest, BizCode: ErrCodeBadRequest.Int()}
}

// IsNotValid is not valid error
func IsNotValid(err error) bool {
	err = Cause(err)
	if innerErr, ok := err.(*CodeError); ok {
		return innerErr.Code == codes.InvalidArgument
	}
	return false
}

// NotFoundf returns an error which satisfies IsNotFound().
func NotFoundf(format string, args ...interface{}) error {
	return &CodeError{Err: wrap(nil, format, " "+http.StatusText(http.StatusNotFound), args...), Code: codes.NotFound, HTTPCode: http.StatusNotFound, BizCode: ErrCodeNotFound.Int()}
}

// NewNotFound returns an error which wraps err that satisfies
func NewNotFound(err error, msg string) error {
	return &CodeError{Err: wrap(err, msg, ""), Code: codes.NotFound, HTTPCode: http.StatusNotFound, BizCode: ErrCodeNotFound.Int()}
}

// IsNotFound is not Fund
func IsNotFound(err error) bool {
	err = Cause(err)
	if innerErr, ok := err.(*CodeError); ok {
		return innerErr.Code == codes.NotFound
	}
	return false
}

// AlreadyExistsf returns an error which satisfies
func AlreadyExistsf(format string, args ...interface{}) error {
	return &CodeError{Err: wrap(nil, format, " "+http.StatusText(http.StatusConflict), args...), Code: codes.AlreadyExists, HTTPCode: http.StatusConflict, BizCode: ErrCodeConflict.Int()}
}

// NewAlreadyExists returns an error which wraps err and satisfies
func NewAlreadyExists(err error, msg string) error {
	return &CodeError{Err: wrap(err, msg, ""), Code: codes.AlreadyExists, HTTPCode: http.StatusConflict, BizCode: ErrCodeConflict.Int()}
}

// IsAlreadyExists is already exists
func IsAlreadyExists(err error) bool {
	err = Cause(err)
	if innerErr, ok := err.(*CodeError); ok {
		return innerErr.Code == codes.AlreadyExists
	}
	return false
}

// Forbiddenf returns an error which satistifes IsForbidden()
func Forbiddenf(format string, args ...interface{}) error {
	return &CodeError{Err: wrap(nil, format, " "+http.StatusText(http.StatusForbidden), args...), Code: codes.PermissionDenied, HTTPCode: http.StatusForbidden, BizCode: ErrCodeForbidden.Int()}
}

// NewForbidden returns an error which wraps err that satisfies
func NewForbidden(err error, msg string) error {
	return &CodeError{Err: wrap(err, msg, ""), Code: codes.PermissionDenied, HTTPCode: http.StatusForbidden, BizCode: ErrCodeForbidden.Int()}
}

// IsForbidden is forbidden error
func IsForbidden(err error) bool {
	err = Cause(err)
	if innerErr, ok := err.(*CodeError); ok {
		return innerErr.Code == codes.PermissionDenied
	}
	return false
}

// FailedPreconditionf returns an error which satisfaction IsFailedPrecondition()
func FailedPreconditionf(format string, args ...interface{}) error {
	return &CodeError{Err: wrap(nil, format, " "+http.StatusText(http.StatusPreconditionFailed), args...), Code: codes.FailedPrecondition, HTTPCode: http.StatusPreconditionFailed, BizCode: ErrCodePreconditionFailed.Int()}
}

// NewFailedPrecondition returns an error which wraps err that satisfies
func NewFailedPrecondition(err error, msg string) error {
	return &CodeError{Err: wrap(err, msg, ""), Code: codes.FailedPrecondition, HTTPCode: http.StatusPreconditionFailed, BizCode: ErrCodePreconditionFailed.Int()}
}

// IsFailedPrecondition is failed precondition errors
func IsFailedPrecondition(err error) bool {
	err = Cause(err)
	if innerErr, ok := err.(*CodeError); ok {
		return innerErr.Code == codes.FailedPrecondition
	}
	return false
}

// Abortedf returns an error which satisfaction IsAborted()
func Abortedf(format string, args ...interface{}) error {
	return &CodeError{Err: wrap(nil, format, " "+http.StatusText(http.StatusInternalServerError), args...), Code: codes.Aborted, HTTPCode: http.StatusInternalServerError, BizCode: ErrCodeInternalServerError.Int()}
}

// NewAborted returns an error which wraps err that satisfies
func NewAborted(err error, msg string) error {
	return &CodeError{Err: wrap(err, msg, ""), Code: codes.Aborted, HTTPCode: http.StatusInternalServerError, BizCode: ErrCodeInternalServerError.Int()}
}

// IsAborted is aborted error
func IsAborted(err error) bool {
	err = Cause(err)
	if innerErr, ok := err.(*CodeError); ok {
		return innerErr.Code == codes.Aborted
	}
	return false
}

// NotImplementedf returns an error which satisfies IsNotImplemented().
func NotImplementedf(format string, args ...interface{}) error {
	return &CodeError{Err: wrap(nil, format, " "+http.StatusText(http.StatusNotImplemented), args...), Code: codes.Unimplemented, HTTPCode: http.StatusNotImplemented, BizCode: ErrCodeNotImplemented.Int()}
}

// NewNotImplemented returns an error which wraps err and satisfies
func NewNotImplemented(err error, msg string) error {
	return &CodeError{Err: wrap(err, msg, ""), Code: codes.Unimplemented, HTTPCode: http.StatusNotImplemented, BizCode: ErrCodeNotImplemented.Int()}
}

// IsNotImplemented is not implemented
func IsNotImplemented(err error) bool {
	err = Cause(err)
	if innerErr, ok := err.(*CodeError); ok {
		return innerErr.Code == codes.Unimplemented
	}
	return false
}

// Internalf returns an error which internal server error
func Internalf(format string, args ...interface{}) error {
	return &CodeError{Err: wrap(nil, format, " "+http.StatusText(http.StatusInternalServerError), args...), Code: codes.Internal, HTTPCode: http.StatusInternalServerError, BizCode: ErrCodeInternalServerError.Int()}
}

// NewInternal == NewInternal return an error which internal server error
func NewInternal(err error, msg string) error {
	if IsBizCodeError(err, MysqlErrorBizCode) { // 对于mysql error, bizcode需要设置为 MysqlErrorBizCode
		return &CodeError{Err: wrap(err, msg, ""), Code: codes.Internal, HTTPCode: http.StatusInternalServerError, BizCode: MysqlErrorBizCode}
	}
	return &CodeError{Err: wrap(err, msg, ""), Code: codes.Internal, HTTPCode: http.StatusInternalServerError, BizCode: ErrCodeInternalServerError.Int()}
}

// IsInternal is internal error
func IsInternal(err error) bool {
	err = Cause(err)
	if innerErr, ok := err.(*CodeError); ok {
		return innerErr.Code == codes.Internal
	}
	return false
}

// Unavailablef returns an error which server unavailable
func Unavailablef(format string, args ...interface{}) error {
	return &CodeError{Err: wrap(nil, format, " "+http.StatusText(http.StatusServiceUnavailable), args...), Code: codes.Unavailable, HTTPCode: http.StatusServiceUnavailable, BizCode: ErrCodeServiceUnavailable.Int()}
}

// NewUnavailable returns an error which server unavailable
func NewUnavailable(err error, msg string) error {
	return &CodeError{Err: wrap(err, msg, ""), Code: codes.Unavailable, HTTPCode: http.StatusServiceUnavailable, BizCode: ErrCodeServiceUnavailable.Int()}
}

// IsUnavailable is unavailable error
func IsUnavailable(err error) bool {
	err = Cause(err)
	if innerErr, ok := err.(*CodeError); ok {
		return innerErr.Code == codes.Unavailable
	}
	return false
}

// Unauthorizedf returns an error which satisfies IsUnauthorized().
func Unauthorizedf(format string, args ...interface{}) error {
	return &CodeError{Err: wrap(nil, format, " "+http.StatusText(http.StatusUnauthorized), args...), Code: codes.Unauthenticated, HTTPCode: http.StatusUnauthorized, BizCode: ErrCodeUnauthorized.Int()}
}

// NewUnauthorized returns an error which wraps err and satisfies
func NewUnauthorized(err error, msg string) error {
	return &CodeError{Err: wrap(err, msg, ""), Code: codes.Unauthenticated, HTTPCode: http.StatusUnauthorized, BizCode: ErrCodeUnauthorized.Int()}
}

// IsUnauthorized is unauthorized
func IsUnauthorized(err error) bool {
	err = Cause(err)
	if innerErr, ok := err.(*CodeError); ok {
		return innerErr.Code == codes.Unauthenticated
	}
	return false
}

// NewBizCodeError new biz code error with biz code and translated message
// Suggest using NewCodeError for more detailed code
func NewBizCodeError(bizCode uint32) error {
	return &CodeError{Err: wrap(nil, "", ""), Code: codes.OK, HTTPCode: http.StatusOK, BizCode: bizCode}
}

// NewBizCodeErrorf new biz code error with biz code and custom message
// Suggest using NewCodeErrorf for more detailed code
func NewBizCodeErrorf(bizCode uint32, format string, args ...interface{}) error {
	return &CodeError{Err: wrap(nil, format, "", args...), Code: codes.OK, HTTPCode: http.StatusOK, BizCode: bizCode}
}

// IsBizCodeError is biz code error
func IsBizCodeError(err error, bizCode uint32) bool {
	err = Cause(err)
	if innerErr, ok := err.(*CodeError); ok {
		return innerErr.BizCode == bizCode
	}
	return false
}

// IsAnyBizCodeErr is any biz code error
func IsAnyBizCodeErr(err error) bool {
	err = Cause(err)
	if innerErr, ok := err.(*CodeError); ok {
		if innerErr.Code == codes.OK && innerErr.HTTPCode == http.StatusOK && innerErr.BizCode > 0 {
			return true
		}
	}
	return false
}
