package main

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"
	wsserver "wstest/proto"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var grpcClientPool sync.Pool
var connCount int

func init() {
	grpcClientPool.New = func() interface{} {
		conn, err := grpc.NewClient("localhost:9992", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		connCount++
		fmt.Printf("创建了一个gRPC连接，目前连接数为：%d\n", connCount)
		return conn
	}
}

func CreateWebSocket(wsid string) {
	u := url.URL{Scheme: "ws", Host: "localhost:9991", Path: "/v1/ws", RawQuery: fmt.Sprintf("wsid=%s", wsid)}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return
	}

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			return
		}
	}
}

var count = 1000

func main() {
	begin := time.Now()

	pool := NewPool(10, 50, Process)

	pool.Start()

	for i := 0; i < count; i++ {
		pool.Submit(fmt.Sprintf("XXXYYYZZZ-%d", i+1))
	}

	pool.Close()

	end := time.Now()
	diff := end.Sub(begin)
	fmt.Printf("消耗秒数为：%f,QPS为：%f", diff.Seconds(), float64(count)/diff.Seconds())

	select {}
}

func Process(msg string) {
	// 实现具体业务逻辑
	conn := grpcClientPool.Get().(*grpc.ClientConn)
	defer grpcClientPool.Put(conn)

	client := wsserver.NewWsServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := client.Preconnect(ctx, &wsserver.PreconnectRequest{Token: msg})
	if err != nil {
		logx.Errorf("调用出错，错误信息为：%v", err)
		return
	}
	fmt.Println(resp.Wsid)

	go CreateWebSocket(resp.Wsid)
}

type Pool struct {
	taskQueue chan string
	taskFn    func(string)
	workers   int
	wg        sync.WaitGroup
}

func NewPool(workers int, capacity int, taskFn func(string)) *Pool {
	pool := Pool{
		taskQueue: make(chan string, 10),
		taskFn:    taskFn,
		workers:   workers,
	}

	pool.wg.Add(workers)
	return &pool
}

func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		go func() {
			defer p.wg.Done()

			for {
				task, ok := <-p.taskQueue
				if !ok {
					return
				}
				p.taskFn(task)
			}
		}()
	}
}

func (p *Pool) Submit(task string) {
	p.taskQueue <- task
}

func (p *Pool) Close() {
	close(p.taskQueue)
	p.wg.Wait()
}
