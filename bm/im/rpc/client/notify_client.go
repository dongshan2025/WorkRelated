package notify

import (
	"context"
	"fmt"
	notify "im/rpc/client/proto"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var notifyPool sync.Pool

func init() {
	notifyPool.New = func() interface{} {
		conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		fmt.Printf("创建了一个gRPC连接: %v\n", conn)
		return conn
	}
}

func GetWsid(token string) (string, error) {
	conn := notifyPool.Get().(*grpc.ClientConn)
	defer notifyPool.Put(conn)

	client := notify.NewNotifyClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.GetWsid(ctx, &notify.GetWsidRequest{Token: "xxxyyyzzz"})
	if err != nil {
		return "", err
	}

	return resp.Wsid, nil
}

func SendNotify(wsids []string, message string) (int32, error) {
	conn := notifyPool.Get().(*grpc.ClientConn)
	defer notifyPool.Put(conn)

	client := notify.NewNotifyClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.SendNotify(ctx, &notify.SendNotifyRequest{Wsid: wsids, Message: message})
	if err != nil {
		return int32(notify.SendNotifyStatus_FAILED), err
	}

	return int32(resp.Status), nil
}
