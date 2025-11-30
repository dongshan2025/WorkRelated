package callback

import (
	"context"
	"fmt"
	"sync"
	"time"
	wscallback "ws/internal/rpc/client/proto"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var callbackPool sync.Pool

func NewCallbackInstance(host string, port int) {
	callbackPool.New = func() any {
		conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logx.Errorf("gRPC Callback客户端连接池初始化失败，错误信息为：%v", err)
			panic(err)
		}

		logx.Infof("gRPC Callback客户端连接池初始化成功，连接信息信息为：%v", conn)
		return conn
	}
}

func Open(wsid string) (int32, error) {
	conn := callbackPool.Get().(*grpc.ClientConn)
	defer callbackPool.Put(conn)
	client := wscallback.NewCallbackClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.Open(ctx, &wscallback.OpenRequest{Wsid: wsid})
	if err != nil {
		return int32(wscallback.CallbackStatus_CALLBACK_FAILED), err
	}

	return int32(resp.Status), nil
}

func Close(wsid string) (int32, error) {
	conn := callbackPool.Get().(*grpc.ClientConn)
	defer callbackPool.Put(conn)
	client := wscallback.NewCallbackClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.Close(ctx, &wscallback.CloseRequest{Wsid: wsid})
	if err != nil {
		return int32(wscallback.CallbackStatus_CALLBACK_FAILED), err
	}

	return int32(resp.Status), nil
}
