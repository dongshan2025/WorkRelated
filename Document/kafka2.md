# 单节点集群配置

# ================================================== KAFKA1 ==================================================
docker run --name=kafka1 \
	-p 19092:9092 -p 19093:9093 \
	-e LANG=C.UTF-8 \
	-e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
	-e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
	-e CLUSTER_ID=82vqfbdSTO2QzS_M0Su1Bw \
	-e KAFKA_NODE_ID=1 \
	-e KAFKA_PROCESS_ROLES=broker,controller \
	-e KAFKA_CONTROLLER_QUORUM_VOTERS="1@192.168.252.104:19093,2@192.168.252.104:29093,3@192.168.252.104:39093" \
	-e KAFKA_LISTENERS="PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093" \
	-e KAFKA_ADVERTISED_LISTENERS="PLAINTEXT://192.168.252.104:19092" \
	apache/kafka:4.0.1============================================= KAFKA2 ==================================================

# ================================================== KAFKA2 ==================================================
docker run --name=kafka2 \
	-p 29092:9092 -p 29093:9093 \
	-e LANG=C.UTF-8 \
	-e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
	-e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
	-e CLUSTER_ID=82vqfbdSTO2QzS_M0Su1Bw \
	-e KAFKA_NODE_ID=2 \
	-e KAFKA_PROCESS_ROLES=broker,controller \
	-e KAFKA_CONTROLLER_QUORUM_VOTERS="1@192.168.252.104:19093,2@192.168.252.104:29093,3@192.168.252.104:39093" \
	-e KAFKA_LISTENERS="PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093" \
	-e KAFKA_ADVERTISED_LISTENERS="PLAINTEXT://192.168.252.104:29092" \
	apache/kafka:4.0.1
# ================================================== KAFKA3 ==================================================
docker run --name=kafka3 \
	-p 39092:9092 -p 39093:9093 \
	-e LANG=C.UTF-8 \
	-e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
	-e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
	-e CLUSTER_ID=82vqfbdSTO2QzS_M0Su1Bw \
	-e KAFKA_NODE_ID=3 \
	-e KAFKA_PROCESS_ROLES=broker,controller \
	-e KAFKA_CONTROLLER_QUORUM_VOTERS="1@192.168.252.104:19093,2@192.168.252.104:29093,3@192.168.252.104:39093" \
	-e KAFKA_LISTENERS="PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093" \
	-e KAFKA_ADVERTISED_LISTENERS="PLAINTEXT://192.168.252.104:39092" \
	apache/kafka:4.0.1
