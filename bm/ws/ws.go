// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"ws/internal/config"
	"ws/internal/global"
	"ws/internal/handler"
	callback "ws/internal/rpc/client"
	wsgrpc "ws/internal/rpc/server"
	"ws/internal/svc"
	"ws/internal/wskafka/wsconsumer"
	"ws/internal/wskafka/wsproducer"
	"ws/pkg/wscache"
	"ws/pkg/wsetcd"
	"ws/pkg/wsuuid"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/ws.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	logx.AddWriter(logx.NewWriter(os.Stdout))

	// 初始化Redis
	cli := wscache.NewWsRedis(c.Redis.Addr, c.Redis.Password, c.Redis.DB)
	pong, err := cli.Ping(context.Background()).Result()
	if err != nil {
		logx.Errorf("连接Redis出现错误，请排查，错误信息为：%v", err)
		return
	}
	logx.Infof("Redis连接成功，pong信息为：%s", pong)

	// 获取WS服务的机器ID
	machineId, err := wscache.GetInstance().Incr(global.WS_MACHINEID)
	if err != nil {
		logx.Errorf("获取WS服务ID失败，错误信息为：%v", err)
		return
	}

	global.MachineID = fmt.Sprintf("%d", machineId)
	global.ApiPort = c.Port
	global.RpcPort = c.RPC.Port
	global.RegisterWsAddr = c.RegisterWsAddr
	global.RegisterRpcIp = c.RegisterRpcIp
	global.KafkaIsOpen = c.Kafka.IsOpen
	global.CallbackIsOpen = c.CallbackRpc.IsOpen

	// 启动WS的gRPC服务
	go wsgrpc.NewRpcWsServer(ctx, &c)

	// 注册服务
	go wsetcd.RegisterService(c.RPC.Name, c.RegisterRpcIp, c.RPC.Port, c.Etcd.Hosts)

	time.Sleep(5 * time.Second)
	logx.Infof("gRPC服务注册成功")

	// 实例化雪花算法
	err = wsuuid.NewUuidInstance(uint16(machineId))
	if err != nil {
		logx.Errorf("实例化雪花算法失败，错误信息为：%v", err)
		return
	}
	logx.Info("雪花算法实例生成成功")

	// 初始化回调服务客户端连接池
	callback.NewCallbackInstance(c.CallbackRpc.Host, c.CallbackRpc.Port)
	logx.Infof("回调RPC客户端连接池初始化成功")

	// Kafka生产者实例化
	wsproducer.NewProducerInstance(c.Kafka.Brokers)
	logx.Infof("Kafka生产者实例化成功")

	// Kafka消费者实例化
	go wsconsumer.NewConsumer(ctx, c.Kafka.Brokers)
	logx.Infof("Kafka消费者实例化成功")

	logx.Infof("WS-API服务在[%s:%d]上启动中...", c.Host, c.Port)

	var wsConnecter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "WebSocket_Count",
		Help: "当前WebSocket连接数",
	})

	global.WebSocketCount = wsConnecter

	server.Start()
}
