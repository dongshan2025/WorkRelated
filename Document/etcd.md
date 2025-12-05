https://etcd.io/docs/v3.4/op-guide/container/
https://etcd.io/docs/v3.7/op-guide/clustering/


===============================================================================
REGISTRY=quay.io/coreos/etcd
# available from v3.6.5

# For each machine
ETCD_VERSION=v3.6.5
TOKEN=my-etcd-token
CLUSTER_STATE=new
NAME_1=etcd-node-0
NAME_2=etcd-node-1
NAME_3=etcd-node-2
HOST_1=192.168.252.101
HOST_2=192.168.252.102
HOST_3=192.168.252.103
CLUSTER=${NAME_1}=http://${HOST_1}:2380,${NAME_2}=http://${HOST_2}:2380,${NAME_3}=http://${HOST_3}:2380
DATA_DIR=/var/lib/etcd

# For node 1
THIS_NAME=${NAME_1}
THIS_IP=${HOST_1}
docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  -d \
  --volume=${DATA_DIR}:/etcd-data \
  --restart always \
  --name etcd ${REGISTRY}:${ETCD_VERSION} \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name ${THIS_NAME} \
  --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster ${CLUSTER} \
  --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
===============================================================================
REGISTRY=quay.io/coreos/etcd
# available from v3.6.5

# For each machine
ETCD_VERSION=v3.6.5
TOKEN=my-etcd-token
CLUSTER_STATE=new
NAME_1=etcd-node-0
NAME_2=etcd-node-1
NAME_3=etcd-node-2
HOST_1=192.168.252.101
HOST_2=192.168.252.102
HOST_3=192.168.252.103
CLUSTER=${NAME_1}=http://${HOST_1}:2380,${NAME_2}=http://${HOST_2}:2380,${NAME_3}=http://${HOST_3}:2380
DATA_DIR=/var/lib/etcd

# For node 2
THIS_NAME=${NAME_2}
THIS_IP=${HOST_2}
docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  -d \
  --volume=${DATA_DIR}:/etcd-data \
  --restart always \
  --name etcd ${REGISTRY}:${ETCD_VERSION} \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name ${THIS_NAME} \
  --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster ${CLUSTER} \
  --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
===============================================================================
REGISTRY=quay.io/coreos/etcd
# available from v3.6.5

# For each machine
ETCD_VERSION=v3.6.5
TOKEN=my-etcd-token
CLUSTER_STATE=new
NAME_1=etcd-node-0
NAME_2=etcd-node-1
NAME_3=etcd-node-2
HOST_1=192.168.252.101
HOST_2=192.168.252.102
HOST_3=192.168.252.103
CLUSTER=${NAME_1}=http://${HOST_1}:2380,${NAME_2}=http://${HOST_2}:2380,${NAME_3}=http://${HOST_3}:2380
DATA_DIR=/var/lib/etcd

# For node 3
THIS_NAME=${NAME_3}
THIS_IP=${HOST_3}
docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  -d \
  --volume=${DATA_DIR}:/etcd-data \
  --restart always \
  --name etcd ${REGISTRY}:${ETCD_VERSION} \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name ${THIS_NAME} \
  --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster ${CLUSTER} \
  --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
===============================================================================
查看集群成员列表：
	docker exec -it etcd etcdctl --write-out=table member list
查看本节点状态
	docker exec -it etcd etcdctl --write-out=table endpoint status
查看本节点健康状态
	docker exec -it etcd etcdctl --write-out=table endpoint health



================================================== Docker Compose方式 ================================================================
https://www.cnblogs.com/crow1840/p/17506992.html


================================================== 基本操作 ==================================================
查看版本：
	docker exec -it etcd etcdctl version

1. 设置值
	docker exec -it etcd etcdctl put /foo/bar "Hello World" #设置/foo/bar的值为“Hello World”
	docker exec -it etcd etcdctl put /foo/bar "Hello World" --ttl 60 #设置/foo/bar的值为“Hello World”，并且超时时间为60s
	docker exec -it etcd etcdctl put /foo/bar "World" --swap-with-value "Hello" #当/foo/bar的值为“Hello”时，则设置/foo/bar的值为“World”，如果不为“Hello”，则会抛出异常

2. 获取值
	docker exec -it etcd etcdctl get /foo/bar
	docker exec -it etcd etcdctl get --prefix /foo #模糊查询匹配到前缀为/foo的数据
	docker exec -it etcd etcdctl --prefix --keys-only=true get /foo #模糊查询匹配到前缀为/foo的键（不返回值）

3. 删除值
	docker exec -it etcd etcdctl del /foo/bar

4. 查看租约列表
	docker exec -it etcd etcdctl lease list
	
5. 创建一个租约
	docker exec -it etcd etcdctl lease grant 120 #创建一个120秒的租约
	lease 225e9a76087e3f10 granted with TTL(120s)

6. 为某个KV授予租约
	docker exec -it etcd etcdctl put web3 "Hello" --lease=225e9a76087e3f10

7. 查看租约信息
	docker exec -it etcd etcdctl lease timetolive 225e9a76087e3f10

8. 重置租约
	docker exec -it etcd etcdctl lease keep-alive 225e9a76087e3f10 #相当于重置了剩余过期时间，所有绑定租约的Key的剩余时间又变成了120秒

9. 撤销租约
	docker exec -it etcd etcdctl lease revoke 225e9a76087e3f10 #租约撤销的同时，被授予租约的KV会被删除

10. 查看所有键值
	docker exec -it etcd etcdctl get --prefix ""








