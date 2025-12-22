# 单机集群配置

# ================================================== ETCD1.SH ==================================================
REGISTRY=quay.io/coreos/etcd
# available from v3.6.5

# For each machine
ETCD_VERSION=v3.6.5
TOKEN=my-etcd-token
CLUSTER_STATE=new
NAME_1=etcd-node-0
NAME_2=etcd-node-1
NAME_3=etcd-node-2
HOST_1=192.168.252.104
HOST_2=192.168.252.104
HOST_3=192.168.252.104
CLUSTER=${NAME_1}=http://${HOST_1}:12380,${NAME_2}=http://${HOST_2}:22380,${NAME_3}=http://${HOST_3}:32380
DATA_DIR=/var/lib/etcd1

# For node 1
THIS_NAME=${NAME_1}
THIS_IP=${HOST_1}
docker run \
  -p 12379:2379 \
  -p 12380:2380 \
  -d \
  --volume=${DATA_DIR}:/etcd-data \
  --restart always \
  --name etcd1 ${REGISTRY}:${ETCD_VERSION} \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name ${THIS_NAME} \
  --initial-advertise-peer-urls http://${THIS_IP}:12380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${THIS_IP}:12379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster ${CLUSTER} \
  --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
# ================================================== ETCD2.SH ==================================================
REGISTRY=quay.io/coreos/etcd
# available from v3.6.5

# For each machine
ETCD_VERSION=v3.6.5
TOKEN=my-etcd-token
CLUSTER_STATE=new
NAME_1=etcd-node-0
NAME_2=etcd-node-1
NAME_3=etcd-node-2
HOST_1=192.168.252.104
HOST_2=192.168.252.104
HOST_3=192.168.252.104
CLUSTER=${NAME_1}=http://${HOST_1}:12380,${NAME_2}=http://${HOST_2}:22380,${NAME_3}=http://${HOST_3}:32380
DATA_DIR=/var/lib/etcd2

# For node 2
THIS_NAME=${NAME_2}
THIS_IP=${HOST_2}
docker run \
  -p 22379:2379 \
  -p 22380:2380 \
  -d \
  --volume=${DATA_DIR}:/etcd-data \
  --restart always \
  --name etcd2 ${REGISTRY}:${ETCD_VERSION} \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name ${THIS_NAME} \
  --initial-advertise-peer-urls http://${THIS_IP}:22380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${THIS_IP}:22379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster ${CLUSTER} \
  --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
# ================================================== ETCD3.SH ==================================================
REGISTRY=quay.io/coreos/etcd
# available from v3.6.5

# For each machine
ETCD_VERSION=v3.6.5
TOKEN=my-etcd-token
CLUSTER_STATE=new
NAME_1=etcd-node-0
NAME_2=etcd-node-1
NAME_3=etcd-node-2
HOST_1=192.168.252.104
HOST_2=192.168.252.104
HOST_3=192.168.252.104
CLUSTER=${NAME_1}=http://${HOST_1}:12380,${NAME_2}=http://${HOST_2}:22380,${NAME_3}=http://${HOST_3}:32380
DATA_DIR=/var/lib/etcd3

# For node 3
THIS_NAME=${NAME_3}
THIS_IP=${HOST_3}
docker run \
  -p 32379:2379 \
  -p 32380:2380 \
  -d \
  --volume=${DATA_DIR}:/etcd-data \
  --restart always \
  --name etcd3 ${REGISTRY}:${ETCD_VERSION} \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name ${THIS_NAME} \
  --initial-advertise-peer-urls http://${THIS_IP}:32380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${THIS_IP}:32379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster ${CLUSTER} \
  --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
