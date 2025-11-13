// go get go.etcd.io/etcd/client/v3
// https://blog.csdn.net/weixin_46501427/article/details/133125020  服务注册与发现
// https://zhuanlan.zhihu.com/p/29278054794

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.252.101:2379", "192.168.252.102:2379", "192.168.252.103:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cli.Close()

	// 创建一个上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动Watch
	watchChan := cli.Watch(ctx, "/config/", clientv3.WithPrefix())
	fmt.Println("开始监听 '/config/' 前缀下的变化...")

	// 处理Watch事件
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Printf("PUT 操作 - Key: %s, Value: %s\n", event.Kv.Key, event.Kv.Value)
			case mvccpb.DELETE:
				fmt.Printf("DELETE 操作 - Key: %s, Value: %s\n", event.Kv.Key, event.Kv.Value)
			}
		}
	}
}
