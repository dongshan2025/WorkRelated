package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func demo() {
	// 1. Establish a connection to etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.252.101:2379", "192.168.252.102:2379", "192.168.252.103:2379"}, // Replace with your etcd endpoint(s)
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer cli.Close()

	// 2. Define the key to watch
	keyToWatch := "/mykey"

	// 3. Start a goroutine to put values into etcd to trigger watch events
	go func() {
		for i := 0; i < 5; i++ {
			value := fmt.Sprintf("value-%d", i)
			_, err := cli.Put(context.Background(), keyToWatch, value)
			if err != nil {
				log.Printf("Error putting key %s: %v", keyToWatch, err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// 4. Create a Watcher
	// Use context.Background() for a watch that persists until the client is closed
	// or use a context with a timeout/cancellation for limited duration watches.
	rch := cli.Watch(context.Background(), keyToWatch)

	// 5. Iterate over the watch channel to receive events
	fmt.Printf("Watching for changes on key: %s\n", keyToWatch)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Event received! Type: %s, Key: %s, Value: %s\n",
				ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}
