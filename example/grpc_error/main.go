package main

import (
	"context"
	"fmt"
	"net"

	"gitee.com/spwx/errors"
	"gitee.com/spwx/errors/grpc_demo"
	"google.golang.org/grpc"
)

var bizErrorCode uint32 = 10086

func main() {
	grpcSrv := grpcServer()
	grpcClient()
	grpcSrv.GracefulStop()
}

type Demo struct {
	grpc_demo.UnimplementedDemoServer
}

func (d Demo) DoDemo(ctx context.Context, req *grpc_demo.Req) (*grpc_demo.Resp, error) {
	// 将error封装成codeError后返回；
	// return nil, errors.ToGRPCReturnError(errors.NewBizCodeErrorf(bizErrorCode, "this is biz code error example"))
	return nil, errors.ToGRPCReturnError(errors.Annotate(errors.AlreadyExistsf(""), "AlreadyExistsf"))
}

// 运行一个grpc服务器
func grpcServer() *grpc.Server {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	grpcSrv := grpc.NewServer()
	grpc_demo.RegisterDemoServer(grpcSrv, &Demo{})
	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return grpcSrv
}

func grpcClient() {
	ctx := context.Background()
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	cli := grpc_demo.NewDemoClient(conn)
	_, err = cli.DoDemo(ctx, &grpc_demo.Req{})
	codeErr := errors.GRPCErrToError(err)

	fmt.Println("errors.IsAlreadyExists(err) ", errors.IsAlreadyExists(codeErr))
	// if errors.IsBizCodeError(codeErr, bizErrorCode) {
	// 	fmt.Println("get the correct biz error code")
	// } else {
	// 	panic("get the wrong biz code")
	// }
}
