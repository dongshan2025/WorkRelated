// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	RPC         RpcConf
	Etcd        EtcdConf
	RegisterIP  string
	Redis       RedisConf
	CallbackRpc CallbackRpcConf
	Kafka       KafkaConf
}

type RpcConf struct {
	Name string
	Host string
	Port int
}

type EtcdConf struct {
	Hosts []string
}

type RedisConf struct {
	Addr     string
	Password string
	DB       int
}

type CallbackRpcConf struct {
	IsOpen bool
	Host   string
	Port   int
}

type KafkaConf struct {
	Brokers []string
	IsOpen  bool
}
