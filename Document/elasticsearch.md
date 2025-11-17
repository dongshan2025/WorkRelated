# 参考资料 https://www.elastic.co/docs/deploy-manage/deploy/self-managed/install-elasticsearch-docker-basic
# 参考资料 https://blog.hanqunfeng.com/2025/03/24/elasticsearch-02-install-cluster/
# 参考资料 https://www.cnblogs.com/dezyan/p/18785095
# 参考资料 https://www.cnblogs.com/chenxingyang/p/17608256.html
# 参考资料 https://www.cnblogs.com/evescn/p/16175547.html
# 参考资料 https://juejin.cn/post/7074115690340286472
# 参考资料 https://blog.csdn.net/sinat_17445041/article/details/149914692
# 参考资料 https://www.cnblogs.com/dengyouf/p/18734266
# 错误处理 https://www.cnblogs.com/wuxizhangjf/p/17771328.html
# Docker安装Elasticsearch三节点集群
## 拉取镜像
docker pull docker.elastic.co/elasticsearch/elasticsearch-wolfi:9.2.1

## 创建网络
docker network create elastic

## 第1个节点配置(主节点)
docker run \
--name es --net elastic -p 9200:9200 -p 9300:9300 -d \
-it -m 1GB \
-e cluster.name=es_cluster \
-e node.name=es01 \
-e network.host=0.0.0.0 \
-e http.port=9200 \
-e transport.port=9300 \
-e network.publish_host=192.168.252.101 \
-v /home/dongshan/es/data:/usr/share/elasticsearch/data \
-v /home/dongshan/es/logs:/usr/share/elasticsearch/logs \
-v /home/dongshan/es/plugins:/usr/share/elasticsearch/plugins \
docker.elastic.co/elasticsearch/elasticsearch-wolfi:9.2.1

正常启动之后，使用以下命令获取集群信息：
docker logs es
集群信息如下：
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ Elasticsearch security features have been automatically configured!
✅ Authentication is enabled and cluster connections are encrypted.

ℹ️  Password for the elastic user (reset with `bin/elasticsearch-reset-password -u elastic`):
  tlSc+ZjJ4Rw578Gr9*bo

ℹ️  HTTP CA certificate SHA-256 fingerprint:
  c48c3dfd3e905df9566d8c88ab4c23c78b5581ca7dc554d9fa2ef10f2417bf55

ℹ️  Configure Kibana to use this cluster:
• Run Kibana and click the configuration link in the terminal when Kibana starts.
• Copy the following enrollment token and paste it into Kibana in your browser (valid for the next 30 minutes):
  eyJ2ZXIiOiI4LjE0LjAiLCJhZHIiOlsiMTkyLjE2OC4yNTIuMTAxOjkyMDAiXSwiZmdyIjoiYzQ4YzNkZmQzZTkwNWRmOTU2NmQ4Yzg4YWI0YzIzYzc4YjU1ODFjYTdkYzU1NGQ5ZmEyZWYxMGYyNDE3YmY1NSIsImtleSI6IkJrZGJoNW9CcDRMSk1xUTctV0FVOm12aFg1WHJJbm1MWkozTmhYMTdIWWcifQ==

ℹ️ Configure other nodes to join this cluster:
• Copy the following enrollment token and start new Elasticsearch nodes with `bin/elasticsearch --enrollment-token <token>` (valid for the next 30 minutes):
  eyJ2ZXIiOiI4LjE0LjAiLCJhZHIiOlsiMTkyLjE2OC4yNTIuMTAxOjkyMDAiXSwiZmdyIjoiYzQ4YzNkZmQzZTkwNWRmOTU2NmQ4Yzg4YWI0YzIzYzc4YjU1ODFjYTdkYzU1NGQ5ZmEyZWYxMGYyNDE3YmY1NSIsImtleSI6IkNFZGJoNW9CcDRMSk1xUTctV0JVOjI0cTEtTjlZRDZET1BJbGh0ZWtLRUEifQ==

  If you're running in Docker, copy the enrollment token and run:
  `docker run -e "ENROLLMENT_TOKEN=<token>" docker.elastic.co/elasticsearch/elasticsearch:9.2.1`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

### 重新生成"elastic"用户的登录密码：
docker exec -it es /usr/share/elasticsearch/bin/elasticsearch-reset-password -u elastic

### 重新生成Kibana进入Token：
docker exec -it es /usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token -s kibana

### 重新生成加入集群的Token：
docker exec -it es /usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token -s node

### 将容器内生成的HTTP证书拷贝到本地：
docker cp es:/usr/share/elasticsearch/config/certs/http_ca.crt .

### 在本地查看集群状态
curl --cacert http_ca.crt -u elastic:tlSc+ZjJ4Rw578Gr9*bo https://localhost:9200/_cat/nodes

## 第2个节点配置
docker run -e ENROLLMENT_TOKEN="eyJ2ZXIiOiI4LjE0LjAiLCJhZHIiOlsiMTkyLjE2OC4yNTIuMTAxOjkyMDAiXSwiZmdyIjoiYzQ4YzNkZmQzZTkwNWRmOTU2NmQ4Yzg4YWI0YzIzYzc4YjU1ODFjYTdkYzU1NGQ5ZmEyZWYxMGYyNDE3YmY1NSIsImtleSI6IkNFZGJoNW9CcDRMSk1xUTctV0JVOjI0cTEtTjlZRDZET1BJbGh0ZWtLRUEifQ==" \
--name es --net elastic -p 9200:9200 -p 9300:9300 -d \
-it -m 1GB \
-e cluster.name=es_cluster \
-e node.name=es02 \
-e network.host=0.0.0.0 \
-e http.port=9200 \
-e transport.port=9300 \
-e network.publish_host=192.168.252.102 \
-v /home/dongshan/es/data:/usr/share/elasticsearch/data \
-v /home/dongshan/es/logs:/usr/share/elasticsearch/logs \
-v /home/dongshan/es/plugins:/usr/share/elasticsearch/plugins \
docker.elastic.co/elasticsearch/elasticsearch-wolfi:9.2.1

## 第3个节点配置
docker run -e ENROLLMENT_TOKEN="eyJ2ZXIiOiI4LjE0LjAiLCJhZHIiOlsiMTkyLjE2OC4yNTIuMTAxOjkyMDAiXSwiZmdyIjoiYzQ4YzNkZmQzZTkwNWRmOTU2NmQ4Yzg4YWI0YzIzYzc4YjU1ODFjYTdkYzU1NGQ5ZmEyZWYxMGYyNDE3YmY1NSIsImtleSI6IkNFZGJoNW9CcDRMSk1xUTctV0JVOjI0cTEtTjlZRDZET1BJbGh0ZWtLRUEifQ==" \
--name es --net elastic -p 9200:9200 -p 9300:9300 -d \
-it -m 1GB \
-e cluster.name=es_cluster \
-e node.name=es03 \
-e network.host=0.0.0.0 \
-e http.port=9200 \
-e transport.port=9300 \
-e network.publish_host=192.168.252.103 \
-v /home/dongshan/es/data:/usr/share/elasticsearch/data \
-v /home/dongshan/es/logs:/usr/share/elasticsearch/logs \
-v /home/dongshan/es/plugins:/usr/share/elasticsearch/plugins \
docker.elastic.co/elasticsearch/elasticsearch-wolfi:9.2.1

## 关闭第一个节点之后，会重新选举新的主节点，这个时候第1个节点要重新加入集群的话需要做以下操作
### 在新的主节点执行以下命令，生成新的集群加入Toeken
docker exec -it es /usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token -s node
### 然后在节点1执行以下命令，注意这个时候"ENROLLMENT_TOKEN"是重新生成的Token
docker run -e ENROLLMENT_TOKEN="eyJ2ZXIiOiI4LjE0LjAiLCJhZHIiOlsiMTkyLjE2OC4yNTIuMTAzOjkyMDAiXSwiZmdyIjoiYmU0NDFhODU3MmZlYjc1ZjBiMmNkOGUzNjExM2FjYWFmYzFhNzQ3NTMwMTViZWYxOTJhZjA0NDI2YzQ5OWZhMiIsImtleSI6IlY0elNocG9CZV93a0FmMU5aV2FxOllNbW1pRnJ4cWhlcU1WamRoZ3NDTFEifQ==" \
-e cluster.name=es_cluster \
-e node.name=es01 \
-e network.host=0.0.0.0 \
-e http.port=9200 \
-e transport.port=9300 \
-e network.publish_host=192.168.252.101 \
--name es --net elastic -p 9200:9200 -p 9300:9300 \
-d --restart always \
-it -m 1GB docker.elastic.co/elasticsearch/elasticsearch-wolfi:9.2.1

## 注意：如果不想丢失数据可以挂载data和log目录
docker run \
-v /home/docker/es/data:/usr/share/elasticsearch/data \
-v /home/docker/es/logs:/usr/share/elasticsearch/logs \
-v /home/docker/es/plugins:/usr/share/elasticsearch/plugins \
## Kibana配置
### 拉取镜像
docker pull docker.elastic.co/kibana/kibana:9.2.1
### 运行容器
docker run --name kibana --net elastic -p 5601:5601 -d --restart always docker.elastic.co/kibana/kibana:9.2.1




