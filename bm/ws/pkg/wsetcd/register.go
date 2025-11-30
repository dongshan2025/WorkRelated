package wsetcd

import (
	"context"
	"fmt"
	"time"
	"ws/internal/global"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func RegisterService(serviceName string, ip string, port int, etcdEndpoints []string) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// 创建一个5秒的租约
	lease, err := client.Grant(context.Background(), 10)
	if err != nil {
		panic(err)
	}

	// 注册服务键值对，将服务名称和地址写入etcd中
	address := fmt.Sprintf("%s:%d", ip, port)
	_, err = client.Put(context.Background(), fmt.Sprintf("%s/%s-%s", serviceName, global.MachineID, address), address, clientv3.WithLease(lease.ID))
	if err != nil {
		panic(err)
	}

	// 定期刷新租约
	leaseCh, err := client.KeepAlive(context.Background(), lease.ID)
	if err != nil {
		panic(err)
	}

	for ch := range leaseCh {
		if ch == nil {
			return
		}
	}
}
