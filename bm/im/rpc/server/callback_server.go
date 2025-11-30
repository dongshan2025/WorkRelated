package callback

import (
	"context"
	"fmt"
	callback "im/rpc/server/proto"
	"net"

	"google.golang.org/grpc"
)

const (
	callbackport = ":40001"
)

type CallbackServer struct {
	callback.UnimplementedCallbackServer
}

func (s *CallbackServer) Open(ctx context.Context, req *callback.OpenRequest) (*callback.OpenResponse, error) {
	fmt.Printf("收到WS端的打开连接通知，wsid为：%s\n", req.Wsid)

	return &callback.OpenResponse{Status: callback.CallbackStatus_SUCCESS}, nil
}
func (s *CallbackServer) Close(ctx context.Context, req *callback.CloseRequest) (*callback.CloseResponse, error) {
	fmt.Printf("收到WS端的关闭连接通知，wsid为：%s\n", req.Wsid)

	return &callback.CloseResponse{Status: callback.CallbackStatus_SUCCESS}, nil
}

func NewCallbackServer() {
	lis, err := net.Listen("tcp", callbackport)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	callback.RegisterCallbackServer(s, &CallbackServer{})

	fmt.Printf("callback server listening at %v\n", lis.Addr())

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
