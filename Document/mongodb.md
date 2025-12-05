docker pull mongodb/mongodb-community-server:8.2-ubi8

docker pull rabbitmq:4.2-rc-alpine

docker pull redis:8.2.2-alpine

docker pull mysql:9.4

docker run --name mongodb -p 27017:27017 -d --restart always mongo:8.2.1

docker run --name mongodb -p 27017:27017 -d --restart always -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=123456 mongodb/mongodb-community-server:8.2-ubi8


docker run --name redis -p 6379:6379 -d --restart always redis:8.2.3

docker run --name mysql -p 3306:3306 -d --restart always -e MYSQL_ROOT_PASSWORD=123456 mysql:9.5.0

docker run --name rabbitmq -p 5672:5672 -p 15672:15672 -d --restart always rabbitmq:4.2-rc-alpine
	进入容器：
		docker exec -it rabbitmq /bin/bash
	启用管理界面：
		rabbitmq-plugins enable rabbitmq_management
	列出用户：
		rabbitmqctl list_users
	添加用户：
		rabbitmqctl add_user admin 123456
	设置用户权限：
		rabbitmqctl set_user_tags admin administrator
		


docker run -d --name couchbase -p 8091-8096:8091-8096 -p 11210-11211:11210-11211 couchbase:enterprise-7.6.7
http://localhost:8091/
Cluster Name: couchbase_local
Administrator / 123456

===========================================================================
Docker部署MongoDB集群

1. 创建MongoDB集群专用网络
	docker network create mongoCluster

docker run --name mongo1 -p 27017:27017 -d --restart always -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=123456 mongodb/mongodb-community-server:8.2-ubi8 mongod --replSet myReplicaSet --bind_ip localhost,mongo1




docker run -d -p 27017:27017 --name mongo1 --network mongoCluster mongodb/mongodb-community-server:8.2-ubi8 mongod --replSet myReplicaSet --bind_ip localhost,mongo1
docker run -d -p 27018:27017 --name mongo2 --network mongoCluster mongodb/mongodb-community-server:8.2-ubi8 mongod --replSet myReplicaSet --bind_ip localhost,mongo2
docker run -d -p 27019:27017 --name mongo3 --network mongoCluster mongodb/mongodb-community-server:8.2-ubi8 mongod --replSet myReplicaSet --bind_ip localhost,mongo3

随便进入一个容器：
	docker exec -it mongo1 bash
	执行命令：
		mongosh
	执行命令：
		rs.initiate({_id: "myReplicaSet", members: [{_id: 0, host: "172.21.32.1:27017"},{_id: 1, host: "172.21.32.1:27018"},{_id: 2, host: "172.21.32.1:27019"}]})
	查看集群状态：
		rs.status()

var uri = "mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=myReplicaSet"
===========================================================================


docker run -d --name prometheus -p 9090:9090 -v D:\docker\prometheu\prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus:latest





docker run --name redis -p 6379:6379 -d --restart always redis:8.2.3
docker run --name mongodb -p 27017:27017 -d --restart always mongo:8.2.1
docker run --name rabbitmq -p 5672:5672 -p 15672:15672 -d --restart always rabbitmq:4.2.0-management

docker run --name es01 -p 9200:9200 -d --restart always -it -m 1GB elasticsearch:9.2.0




curl --cacert http_ca.crt -u elastic:rw9mdcG=3a_spxspj+OR https://localhost:9200


docker run --name kib01 -p 5601:5601 -d --restart always kibana:9.2.0

Password for the [elastic] user successfully reset.
New value: rw9mdcG=3a_spxspj+OR
root@ubuntu:/home/dongshan/es01# docker exec -it es01 /usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token -s kibana
eyJ2ZXIiOiI4LjE0LjAiLCJhZHIiOlsiMTcyLjE3LjAuMzo5MjAwIl0sImZnciI6ImI4NjQ1MmE5YzgyZDAxOWE4ODM3NzI4Y2RkNzFjM2E4MDQ5ZjQyNzhkNzBjOWRhZjk2YTIyYjVmZjU4MWJiMzYiLCJrZXkiOiJiUk5YZTVvQnBhVmFGQms0MUpsYzpOc0lYNnJQREoyUVc4MkFPeVRyYWVBIn0=
root@ubuntu:/home/dongshan/es01# docker cp es01:/usr/share/elasticsearch/config/certs/http_ca.crt .
Successfully copied 3.58kB to /home/dongshan/es01/.
root@ubuntu:/home/dongshan/es01# ls
http_ca.crt
root@ubuntu:/home/dongshan/es01# curl --cacert http_ca.crt -u elastic:rw9mdcG=3a_spxspj+OR https://localhost:9200
{
  "name" : "e80171af3462",
  "cluster_name" : "docker-cluster",
  "cluster_uuid" : "ZcmKutQ1TZywzZGY-EQn_Q",
  "version" : {
    "number" : "9.2.0",
    "build_flavor" : "default",
    "build_type" : "docker",
    "build_hash" : "25d88452371273dd27356c98598287b669a03eae",
    "build_date" : "2025-10-21T10:06:21.288851013Z",
    "build_snapshot" : false,
    "lucene_version" : "10.3.1",
    "minimum_wire_compatibility_version" : "8.19.0",
    "minimum_index_compatibility_version" : "8.0.0"
  },
  "tagline" : "You Know, for Search"
}


docker run -e ENROLLMENT_TOKEN="eyJ2ZXIiOiI4LjE0LjAiLCJhZHIiOlsiMTcyLjE3LjAuMzo5MjAwIl0sImZnciI6ImI4NjQ1MmE5YzgyZDAxOWE4ODM3NzI4Y2RkNzFjM2E4MDQ5ZjQyNzhkNzBjOWRhZjk2YTIyYjVmZjU4MWJiMzYiLCJrZXkiOiJvQk5vZTVvQnBhVmFGQms0ekpubTpwSXFaOEszUEhYa1N6cl9YMVJQeG5BIn0=" --name es02 -p 9200:9200 -d --restart always -it -m 1GB elasticsearch:9.2.0


curl --cacert http_ca.crt -u elastic:rw9mdcG=3a_spxspj+OR https://localhost:9200/_cat/nodes

curl -k -u elastic:rw9mdcG=3a_spxspj+OR https://localhost:9200/_cat/nodes

curl --cacert http_ca.crt -u elastic:rw9mdcG=3a_spxspj+OR https://localhost:9200/_cluster/health?pretty


docker run --name es01 -p 9991:9991 -p 9992:9992 -v /home/dongshan/ws/etc/ws.yaml:/app/etc/ws.yaml ws:v1











