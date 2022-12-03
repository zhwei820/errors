package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var grpcCodeToHttpCode = map[codes.Code]uint32{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           http.StatusRequestTimeout,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusRequestTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.ResourceExhausted:  http.StatusGone,
	codes.FailedPrecondition: http.StatusPreconditionFailed,
	codes.Aborted:            http.StatusForbidden,
	codes.OutOfRange:         http.StatusPreconditionFailed,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusGone,
	codes.Unauthenticated:    http.StatusUnauthorized,
}

// ToGRPCStatus to grpc status error
func ToGRPCStatus(err error) *status.Status {
	err = Cause(err)
	if err == nil {
		return nil
	}
	inner, ok := err.(*CodeError)
	if ok {
		if inner.Code == codes.OK && inner.BizCode != 0 {
			inner.Code = codes.Unknown
		}
		st := status.New(inner.Code, err.Error())
		if inner.BizCode != 0 {
			st, _ = st.WithDetails(&BizErrorCode{Code: inner.BizCode})
		}
		return st
	}
	st, _ := status.FromError(err)
	return st
}

// ToGRPCReturnError generate grpc api return error
func ToGRPCReturnError(err error) error {
	st := ToGRPCStatus(err)
	if st == nil {
		return nil
	}
	return st.Err()
}

// GRPCErrToError grpc error to code error
func GRPCErrToError(err error) error {
	if err == nil {
		return nil
	}
	st, _ := status.FromError(err)
	httpCode := grpcCodeToHttpCode[st.Code()]
	codeErr := &CodeError{Err: wrap(nil, err.Error(), ""), Code: st.Code(), HTTPCode: httpCode}
	details := st.Details()
	for _, detail := range details {
		if bizErrorCode, ok := detail.(*BizErrorCode); ok {
			codeErr.BizCode = bizErrorCode.Code
		}
	}
	return codeErr
}
