package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	ServiceName = "/PushService"
)

type ServiceDiscovery struct {
	client     *clientv3.Client
	lock       sync.Mutex
	services   map[string]string
	watchedKey string
}

// 创建服务对象
func NewServiceDiscovery(endpoints []string, watchedKey string) (*ServiceDiscovery, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		return nil, err
	}

	return &ServiceDiscovery{
		client:     client,
		services:   make(map[string]string),
		watchedKey: watchedKey,
	}, nil
}

// 服务发现
func (sd *ServiceDiscovery) DiscoveryService() error {
	resp, err := sd.client.Get(context.Background(), ServiceName, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	sd.lock.Lock()
	for _, kv := range resp.Kvs {
		sd.services[string(kv.Key)] = string(kv.Value)
	}
	defer sd.lock.Unlock()

	return nil
}

// 监听
func (sd *ServiceDiscovery) WatchServices() error {
	watchCh := sd.client.Watch(context.Background(), sd.watchedKey, clientv3.WithPrefix())
	for watchResp := range watchCh {
		for _, event := range watchResp.Events {
			switch event.Type {
			case clientv3.EventTypePut:
				sd.handlePutEvent(event)
			case clientv3.EventTypeDelete:
				sd.handleDeleteEvent(event)
			}
		}
	}

	return nil
}

// 新增或者修改服务
func (sd *ServiceDiscovery) handlePutEvent(event *clientv3.Event) {
	sd.lock.Lock()
	defer sd.lock.Unlock()
	serviceName := string(event.Kv.Key)
	serviceAddress := string(event.Kv.Value)
	sd.services[serviceName] = serviceAddress
	fmt.Printf("Added service: %s -> %s\n", serviceName, serviceAddress)
}

// 删除服务
func (sd *ServiceDiscovery) handleDeleteEvent(event *clientv3.Event) {
	sd.lock.Lock()
	defer sd.lock.Unlock()
	serviceName := string(event.Kv.Key)
	delete(sd.services, serviceName)
	fmt.Printf("Removed service: %s\n", serviceName)
}

func main() {
	serviceDiscovery, err := NewServiceDiscovery([]string{"192.168.252.101:2379", "192.168.252.102:2379", "192.168.252.103:2379"}, ServiceName)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = serviceDiscovery.DiscoveryService()
	if err != nil {
		log.Fatal(err)
		return
	}
	go func() {
		err := serviceDiscovery.WatchServices()
		if err != nil {
			fmt.Printf("Failed to watch services: %v", err)
		}
	}()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	for {
		select {
		case <-ticker.C:
			fmt.Println("服务最新列表：")
			for _, v := range serviceDiscovery.services {
				fmt.Println(v)
			}
		case <-sig:
			fmt.Println("Program terminated.")
			return
		}
	}
}
