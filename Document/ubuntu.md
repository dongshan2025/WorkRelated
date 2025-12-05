
1. 配置root用户
	$ sudo passwd root
	New password:
	Retype new password:
	passwd: password update successfully
	$ su
	Password:

2. 配置ssh
	我们安装系统的时候勾选了要安装SSH，所以系统自动安装了SSH
	2.1 查看版本
		ssh -V
	2.2 查看服务状态
		sudo systemctl status ssh
	2.3 开机启动服务
		sudo systemctl enable ssh
	2.4 关闭开机启动服务
		sudo systemctl disable ssh
	2.5 启动服务
		sudo systemctl start ssh
	2.6 重启服务
		sudo systemctl restart ssh
	2.7 停止服务
		sudo systemctl stop ssh

3. 配置防火墙
	3.1 查看防火墙状态
		sudo ufw status
		如果显示“sudo ufw status”说明没有启用
	3.2 启用防火墙
		sudo ufw enable
	3.3 关闭防火墙
		sudo ufw disable
	3.4 允许SSH连接(开放端口22)
		sudo ufw allow ssh or sudo ufw allow 22/tcp
	3.5 查看当前所有规则
		sudo ufw status verbose
	3.6 允许特定端口访问
		sudo ufw allow 80/tcp
	3.7 删除访问规则
		sudo ufw delete allow ssh
		sudo ufw delete allow 80/tcp
	3.8 重新加载规则
		sudo ufw reload
	3.9 添加带IP限制的规则
		sudo ufw allow from 192.168.252.132 to any port 22
		sudo ufw allow from 192.168.2.0/24 to any port 22 proto tcp
		sudo ufw allow from 10.105.28.0/24 to 10.105.29.1 port 22 proto tcp comment 'Allow SSH'
	3.10 删除带IP限制的规则
		sudo ufw status numbered
		sudo ufw delete 数字
		sudo ufw delete allow from 192.168.252.132

4. 安装net-tools
	sudo apt install net-tools
	
5. 配置静态IP地址
	sudo ls /etc/netplan/*.yaml
	sudo nano /etc/netplan/50-cloud-init.yaml
	修改文件如下：
		network:
		  version: 2
		  renderer: networkd
		  ethernets:
			ens33:
			  dhcp4: no
			  addresses:
				- 192.168.252.100/24
			  routes:
				- to: default
				  via: 192.168.252.2
			  nameservers:
				addresses:
				  - 192.168.252.2
	应用配置：
	sudo netplan apply
192.168.26.128

	查看默认网关：
	sudo ip route show default
	
	查看DNS
	sudo resolvectl status

	退出nano编辑器命令：Ctrl+X

6. 安装Docker
	6.1 Set up Docker's apt repository.
		# Add Docker's official GPG key:
		sudo apt-get update
		sudo apt-get install ca-certificates curl
		sudo install -m 0755 -d /etc/apt/keyrings
		sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
		sudo chmod a+r /etc/apt/keyrings/docker.asc

		# Add the repository to Apt sources:
		echo \
		  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
		  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
		  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
		
		sudo apt-get update
	6.2 Install the Docker packages.
		sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
	6.3 To verify that Docker is running
		sudo systemctl status docker
		sudo systemctl start docker
	6.4 Verify that the installation is successful by running the hello-world image
		sudo docker run hello-world

	6.5 解决国内拉取不到镜像问题
		sudo mkdir -p /etc/docker
		sudo vim /etc/docker/daemon.json
		在daemon.json文件中输入如下内容：
{
  "registry-mirrors": [
	"http://hub-mirror.c.163.com",
	"https://mirrors.tuna.tsinghua.edu.cn",
	"http://mirrors.sohu.com",
	"https://ustc-edu-cn.mirror.aliyuncs.com",
	"https://ccr.ccs.tencentyun.com",
	"https://docker.m.daocloud.io",
	"https://docker.awsl9527.cn"
  ]
}

	重载生效：
		sudo systemctl daemon-reload
	重启服务：
		sudo systemctl restart docker
	查看是否配置成功：
		sudo docker info

================================================== 常用镜像 ==================================================
docker pull jenkins/jenkins:jdk21										https://hub.docker.com/r/jenkins/jenkinsdoc
docker pull mysql:9.5.0													https://hub.docker.com/_/mysql
docker pull docker.elastic.co/elasticsearch/elasticsearch-wolfi:9.2.1	https://www.elastic.co/docs/deploy-manage/deploy/self-managed/install-elasticsearch-docker-basic
docker pull docker.elastic.co/kibana/kibana:9.2.1						https://www.elastic.co/docs/deploy-manage/deploy/self-managed/install-elasticsearch-docker-basic
docker pull apache/kafka:4.0.1											https://hub.docker.com/r/apache/kafka
docker pull prom/prometheus:v3.7.3										https://hub.docker.com/r/prom/prometheus
docker pull mongo:8.2.1													https://hub.docker.com/_/mongo
docker pull redis:8.2.3													https://hub.docker.com/_/redis
docker pull rabbitmq:4.2.0-management									https://hub.docker.com/_/rabbitmq
docker pull quay.io/coreos/etcd:v3.6.5									https://quay.io/repository/coreos/etcd?tab=tags&tag=latest
docker pull memcached:1.6.39											https://hub.docker.com/_/memcached
docker pull nginx:1.29.3-perl											https://hub.docker.com/_/nginx
docker pull jaegertracing/all-in-one:1.74.0								https://hub.docker.com/r/jaegertracing/all-in-one

docker pull prom/alertmanager:v0.29.0
docker pull prom/prometheus:v3.7.3
docker pull grafana/grafana:12.4.0-19363970803
docker pull prom/node-exporter:v1.10.2

方法1：
# 导出镜像到tar文件
docker save -o rabbitmq.tar rabbitmq:4.2.0
# 将tar文件scp到另一台机器
scp rabbitmq.tar dongshan@192.168.252.101:/home/dongshan
# 从tar文件导入镜像到Docker
docker load < rabbitmq.tar
或 docker load -i /home/dongshan/rabbitmq.tar

方法2：
# 导出镜像到tar文件
docker export rabbitmq:4.2.0 > rabbitmq.tar
# 导入镜像
docker import /home/dongshan/rabbitmq.tar
或 cat rabbitmq.tar | docker import - rabbitmq:4.2.0(镜像名称自己定义)
================================================== 磁盘扩容 ==================================================
扩容 /dev/mapper/ubuntu--vg-ubuntu--lv

1. 查看磁盘情况
	df -h
2. 查看lvm卷组的信息
	vgdisplay
		--- Volume group ---
		  VG Name               ubuntu-vg
		  System ID             
		  Format                lvm2
		  Metadata Areas        1
		  Metadata Sequence No  3
		  VG Access             read/write
		  VG Status             resizable
		  MAX LV                0
		  Cur LV                1
		  Open LV               1
		  Max PV                0
		  Cur PV                1
		  Act PV                1
		  VG Size               <38.00 GiB
		  PE Size               4.00 MiB
		  Total PE              9727
		  Alloc PE / Size       7423 / <29.00 GiB
		  Free  PE / Size       2304 / 9.00 GiB
		  VG UUID               3O1Pi0-cdRh-XReq-e886-hs08-JD0k-fDZRc4
	Free  PE / Size       2304 / 9.00 GiB 就是还可以扩充的大小
3. 使用命令进行扩容
	lvextend -L +10G /dev/mapper/ubuntu--vg-ubuntu--lv // 增加10G
	lvextend -L -10G /dev/mapper/ubuntu--vg-ubuntu--lv // 减少10G
4. 执行扩容命令
	resize2fs /dev/mapper/ubuntu--vg-ubuntu--lv
5. 再次查看lv卷组信息
	vgdisplay




























