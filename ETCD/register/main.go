package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	ServiceName = "/PushService"
)

// 获取本机空闲端口号，模拟多服务
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	cli, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer cli.Close()
	return cli.Addr().(*net.TCPAddr).Port, nil
}

func RegisterService(client *clientv3.Client, address string) error {
	// 创建一个租约2秒
	resp, err := client.Grant(context.Background(), 2)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// 注册服务键值对，将服务名称和地址写入etcd中
	_, err = client.Put(context.Background(), fmt.Sprintf("%s/%s", ServiceName, address), address, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
		return err
	}

	// 定期刷新租约
	ch, err := client.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	go func() {
		for resp := range ch {
			if resp == nil {
				log.Println("KeepAlive channel closed")
				return
			}
		}
	}()

	return nil
}

func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.252.101:2379", "192.168.252.102:2379", "192.168.252.103:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Close()

	port, err := GetFreePort()
	if err != nil {
		log.Fatal(err)
		return
	}

	// 注册服务
	err = RegisterService(client, fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("服务注册成功：localhost:%d", port)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	fmt.Println("Program terminated.")
}
