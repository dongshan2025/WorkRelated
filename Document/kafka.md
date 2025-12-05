# Docker安装Kafka集群(三节点)

## 节点1
docker run --name=kafka \
	-p 9092:9092 -p 9093:9093 \
	-e LANG=C.UTF-8 \
	-e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
	-e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
	-e CLUSTER_ID=82vqfbdSTO2QzS_M0Su1Bw \
	-e KAFKA_NODE_ID=1 \
	-e KAFKA_PROCESS_ROLES=broker,controller \
	-e KAFKA_CONTROLLER_QUORUM_VOTERS="1@192.168.252.101:9093,2@192.168.252.102:9093,3@192.168.252.103:9093" \
	-e KAFKA_LISTENERS="PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093" \
	-e KAFKA_ADVERTISED_LISTENERS="PLAINTEXT://192.168.252.101:9092" \
	apache/kafka:4.0.1

## 节点2
docker run --name=kafka \
	-p 9092:9092 -p 9093:9093 \
	-e LANG=C.UTF-8 \
	-e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
	-e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
	-e CLUSTER_ID=82vqfbdSTO2QzS_M0Su1Bw \
	-e KAFKA_NODE_ID=2 \
	-e KAFKA_PROCESS_ROLES=broker,controller \
	-e KAFKA_CONTROLLER_QUORUM_VOTERS="1@192.168.252.101:9093,2@192.168.252.102:9093,3@192.168.252.103:9093" \
	-e KAFKA_LISTENERS="PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093" \
	-e KAFKA_ADVERTISED_LISTENERS="PLAINTEXT://192.168.252.102:9092" \
	apache/kafka:4.0.1

## 节点3
docker run --name=kafka \
	-p 9092:9092 -p 9093:9093 \
	-e LANG=C.UTF-8 \
	-e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
	-e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
	-e CLUSTER_ID=82vqfbdSTO2QzS_M0Su1Bw \
	-e KAFKA_NODE_ID=3 \
	-e KAFKA_PROCESS_ROLES=broker,controller \
	-e KAFKA_CONTROLLER_QUORUM_VOTERS="1@192.168.252.101:9093,2@192.168.252.102:9093,3@192.168.252.103:9093" \
	-e KAFKA_LISTENERS="PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093" \
	-e KAFKA_ADVERTISED_LISTENERS="PLAINTEXT://192.168.252.103:9092" \
	apache/kafka:4.0.1

# 常用命令
在任何一个节点执行以下命令进入Docker：
docker exec -it kafka bash
cd /opt/kafka/bin

## 查看Kafka版本
./kafka-run-class.sh --version

## 创建Topic
./kafka-topics.sh --bootstrap-server localhost:9092 --create --topic my-topic --partitions 3 --replication-factor 3 // 创建3个分区，3个副本的Topic
./kafka-topics.sh --bootstrap-server localhost:9092 --create --topic my-topic --partitions 3 --replication-factor 3 --config max.message.bytes=64000 // 同时指定配置
./kafka-topics.sh --bootstrap-server localhost:9092 --create --topic my-topic --partitions 3 --replication-factor 3 --config x=y // 指定自定义配置
./kafka-topics.sh --bootstrap-server localhost:9092 --alter --topic my-topic --partitions 5 // 将分区扩大到5个

## 查看Topic列表
./kafka-topics.sh --bootstrap-server localhost:9092 --list

## 查看Topic详情
./kafka-topics.sh --bootstrap-server localhost:9092 --describe --topic my-topic

## 删除Topic
./kafka-topics.sh --bootstrap-server localhost:9092 --delete --topic my-topic

## 修改Topic
./kafka-configs.sh --bootstrap-server localhost:9092 --entity-type topics --entity-name my-topic --alter --add-config max.message.bytes=128000

## 删除配置
./kafka-configs.sh --bootstrap-server localhost:9092 --entity-type topics --entity-name my-topic --alter --delete-config max.message.bytes

## 查看修改后的详情
./kafka-configs.sh --bootstrap-server localhost:9092 --entity-type topics --entity-name my-topic --describe

## 添加配置
./kafka-configs.sh --bootstrap-server localhost:9092 --entity-type topics --entity-name my-topic --alter --add-config x=y

## 移除配置
./kafka-configs.sh --bootstrap-server localhost:9092 --entity-type topics --entity-name my-topic --alter --delete-config x

## 创建生产者
./kafka-console-producer.sh --bootstrap-server localhost:9092 --topic my-topic
./kafka-console-producer.sh --bootstrap-server localhost:9092 --topic my-topic --group console-consumer-32439 // 指定消费者组

## 创建消费者
./kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic my-topic // 从当前消息开始消费
./kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic my-topic --from-beginning // 从头开始消费

./kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic my-topic --offset latest // 从尾部开始消费
./kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic my-topic --offset latest --partition 1 // 指定分区
./kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic my-topic --offset latest --partition 1 --max-messages 1 // 只取1条信息
./kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic my-topic --from-beginning --group console-consumer-32439 // 指定消费者组


## 查看消费者组列表
./kafka-groups.sh --bootstrap-server localhost:9092 --list

## 查看某个消费者组详情
./kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group my-group // 查看消费者组消费详情
./kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group my-group --members // 查看消费者组成员
./kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group my-group --members --verbose
./kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group my-group --state
## 删除某个消费者组
./kafka-consumer-groups.sh --bootstrap-server localhost:9092 --delete --group my-group --group my-other-group

## 重置offset
./kafka-consumer-groups.sh --bootstrap-server localhost:9092 --reset-offsets --group my-group --topic my-topic --to-latest

## 查看副本信息
./kafka-metadata-quorum.sh --bootstrap-server localhost:9092 describe --replication


## 查看集群ID
./kafka-cluster.sh cluster-id --bootstrap-server localhost:9092










docker run -p 9092:9092 -d --restart always --name kafka apache/kafka:4.0.1
