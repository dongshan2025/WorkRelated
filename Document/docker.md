# ************************************************** 镜像仓库 **************************************************
# 拉取镜像
docker pull mysql:9.4

# 登录Docker Hub
docker login --username myusername

# 生成Docker Image
docker build -t dongshan2025/dongshan:v1 .

# 推送镜像到Docker Hub
docker push dongshan2025/dongshan:v1

# 搜索镜像
docker search dongshan2025/dongshan:v1

# 登出Docker Hub
docker logout
# ************************************************** 本地镜像 **************************************************
# 列出本地所有镜像
docker images

# 列出本地所有镜像，并带有摘要信息
docker images --digests

# 列出本地所有镜像，包括中间层镜像
docker images -a 或 docker images --all

# 列出本地所有镜像，只显示镜像ID
docker images -q 或 docker images --quiet

# 列出本地所有镜像，并显示完整的镜像ID
docker images --no-trunc

# 根据过滤条件列出本地所有镜像
docker images --filter "reference=redis"        // 根据镜像的名称或标签进行过滤
docker images --filter "since=nginx:1.21"       // 查找在nginx:1.21之后构建的所有镜像
docker images --filter "before=ubuntu:18.04"    // 查找在ubuntu:18.04之前构建的所有镜像
docker images --filter "label=maintainer=admin" // 查找所有maintainer标签值为admin的镜像

# 使用自定义格式列出本地所有镜像
docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.Size}}"

# 删除单个镜像
docker rmi ubuntu:latest

# 删除多个镜像
docker rmi ubuntu:latest nginx:latest

# 删除指定仓库名下的所有镜像
docker rmi -a ubuntu 或 docker rmi --all-tags ubuntu

# 强制删除镜像，即使该镜像正在被容器使用
docker rmi -f ubuntu:latest 或 docker rmi --force ubuntu:latest

# 删除没有标签的悬空镜像
docker rmi -d

# 删除镜像但是保留其子镜像
docker rmi --no-prune ubuntu:latest

# 为镜像打标签
docker tag myimage:1.0 myrepo/myimage:latest

# 为镜像打标签以推送到Docker Hub
docker tag my-app:v1 dongshan2025/dongshan:latest
docker push dongshan2025/dongshan:latest

# 构建镜像
docker build -t myimage:latest .

# 指定Dockerfile文件构建镜像
docker build -f /path/to/Dockerfile -t myimage:latest .

# 设置构建参数构建镜像
docker build --build-arg HTTP_PROXY=http://proxy.example.com -t myimage:latest . // 这会在构建过程中使用 HTTP_PROXY 环境变量

# 不使用缓存层构建镜像
docker build --no-cache -t myimage:latest . // 这会在构建镜像时忽略所有缓存层，确保每一步都重新执行

# 查看镜像历史
docker history myimage:latest

# 查看镜像历史，并显示完整输出
docker history --no-trunc myimage:latest

# 查看镜像历史，仅显示镜像ID
docker history -q myimage:latest

# 保存单个镜像到文件
docker save -o myimage.tar myimage:latest

# 保存多个镜像到文件
docker save -o multiple_images.tar myimage:latest anotherimage:latest

# 从文件加载镜像
docker load -i myimage.tar

# 从标准输入加载镜像
cat myimage.tar | docker load

# 导入镜像，
# docker import命令用于从一个tar文件或URL导入容器快照，从而创建一个新的Docker镜像
# 与docker load不同，docker import可以从容器快照中创建新的镜像，而不需要保留镜像的历史和元数据
docker import mycontainer.tar mynewimage:latest // 这将从 mycontainer.tar 文件导入镜像，并命名为 mynewimage:latest

# 从URL导入镜像
docker import http://example.com/mycontainer.tar mynewimage:latest // 这将从指定的 URL 导入镜像，并命名为 mynewimage:latest

# 从标准输入导入镜像
cat mycontainer.tar | docker import - mynewimage:latest // 这将通过管道从标准输入读取 tar 文件并导入镜像

# 在导入过程中应用变更
docker import -c "ENV LANG=en_US.UTF-8" -c "CMD /bin/bash" mycontainer.tar mynewimage:latest // 这将从 mycontainer.tar 导入镜像，并在过程中设置环境变量 LANG 和命令 CMD

# ---------- 实列开始 ----------
# 创建并运行一个容器
docker run -d --name mycontainer ubuntu:20.04 sleep 3600
# 导出容器快照
docker export mycontainer -o mycontainer.tar
# 从容器快照导入以创建新的镜像
docker import mycontainer.tar mynewimage:latest
# 验证导入的镜像
docker images
# 运行导入的镜像
docker run -it mynewimage:latest /bin/bash
# ---------- 实列结束 ----------

# ************************************************** 容器周期 **************************************************
# 创建并启动一个新容器
docker run --name mycontainer -d -p 8080:8080 --network mynetwork --restart always my-app:v1

# 容器停止后自动删除容器
docker run --name mycontainer -d -p 8080:8080 --rm my-app:v1

# 挂载卷
docker run --name mycontainer -d -p 8080:8080 -v /host/data:/container/data my-app:v1

# 设置环境变量
docker run -e MY_ENV_VAR=my_value --name mycontainer -d -p 8080:8080 my-app:v1

# 交互式运行并分配终端
docker run --name mycontainer -p 8080:8080 -it my-app:v1 /bin/bash

# 启动一个已创建的容器
docker start mycontainer

# 启动多个已创建的容器
docker start mycontainer mycontainer2

# 停止一个正在运行的容器
docker stop mycontainer

# 等待指定时间后停止容器
docker stop -t 30 mycontainer // 等待30秒后停止容器

# 停止多个容器
docker stop mycontainer mycontainer2

# 重启一个容器
docker restart mycontainer

# 重启多个容器
docker restart mycontainer mycontainer2

# 等待指定时间后重启容器
docker restart -t 30 mycontainer // 等待30秒后重启容器

# 立即终止容器的运行
docker kill mycontainer mycontainer2

# 发送信号终止容器的运行
docker kill -s SIGTERM mycontainer
    # SIGKILL: 强制终止进程（默认信号）。
    # SIGTERM: 请求进程终止。
    # SIGINT: 发送中断信号，通常表示用户请求终止。
    #SIGHUP: 挂起信号，通常表示终端断开。

# 删除已经停止的容器
docker rm mycontainer mycontainer2

# 删除所有已停止的容器
docker container prune 或 docker rm $(docker ps -a -q)

# 强制删除一个正在运行的容器
docker rm -f mycontainer // 使用SIGKILL信号

# 暂停一个或多个容器
docker pause mycontainer mycontainer2

# 恢复一个或多个容器
docker unpause mycontainer mycontainer2

# 创建一个容器但不启动
docker create --name mycontainer -d -p 8080:8080 --restart always myapp:v1

# 在运行的容器中执行一个新命令
docker exec -it mycontainer /bin/bash // 以交互模式运行命令
docker exec mycontainer ls /app // 在运行的mycontainer容器中执行ls /app命令，列出/app目录的内容
docker exec -d mycontainer touch /app/newfile.txt // 在运行中的 my_container 容器内后台执行 touch /app/newfile.txt 命令，创建一个新文件
docker exec -e MY_ENV_VAR=my_value my_container env // 在运行中的 my_container 容器内执行 env 命令，并设置环境变量 MY_ENV_VAR 的值为 my_value
docker exec -w /app my_container pwd // 在运行中的 my_container 容器内以 /app 目录为工作目录执行 pwd 命令

# 重命名容器
docker rename my_old_container my_new_container // 根据容器名称重命名
docker rename 123abc456def my_new_container // 根据容器ID重命名
# ************************************************** 容器操作 **************************************************
# 列出所有正在运行的容器
docker ps

# 列出所有容器，包括已经停止的容器
docker ps -a

# 列出所有容器，但只显示容器ID
docker ps -a -q

# 列出最近创建的n个容器
docker ps -n 3

# 容器的7种状态
    created（已创建）
    restarting（重启中）
    running（运行中）
    removing（迁移中）
    paused（暂停）
    exited（停止）
    dead（死亡）

# 获取Docker对象（容器、镜像、卷、网络等）的详细信息
# 检查容器
docker inspect mycontainer

# 检查镜像
docker inspect myimage

# 检查卷
docker inspect myvolume

# 检查网络
docker inspect mynetwork

# 显示指定容器中的正在运行的进程
docker top mycontainer
docker top mycontainer -o pid,cmd // 只显示pid和cmd列

# 实时获取Docker守护进程生成的事件
docker events
docker events --filter event=stop // 过滤事件
docker events --since "2023-07-22T15:04:05" // 显示从 2023-07-22T15:04:05 开始的事件
docker events --until "2023-07-22T16:04:05" // 显示直到 2023-07-22T16:04:05 的事件

# 查看容器的日志输出
docker logs mycontainer
docker logs -f mycontainer // 跟随日志输出
docker logs -t mycontainer // 显示包含时间戳的日志
docker logs --since="2023-07-22T15:00:00" mycontainer // 显示 2023-07-22T15:00:00 之后的日志
docker logs --tail 10 mycontainer // 显示最后10行日志
docker logs --details mycontainer // 包含额外详细信息
docker logs --until="2023-07-22T16:00:00" mycontainer // 显示 2023-07-22T16:00:00 之前的日志

# 等待一个容器停止并获取其退出代码
docker wait mycontainer
# docker wait 命令用于阻塞，直到指定的容器停止运行，然后返回容器的退出代码。
# docker wait 命令对于自动化脚本非常有用，因为它可以等待容器完成某项任务，并根据容器的退出状态采取后续操作。
# docker wait 命令是一个简单但非常有用的工具，允许用户等待容器停止并获取其退出代码。通过该命令，用户可以轻松地在脚本中实现任务同步和自动化操作。使用 docker wait 命令，可以确保在指定的容器完成其任务之前，不会进行任何后续操作。
docker run --name test_container ubuntu bash -c "exit 5" // 启动一个会立即退出的容器
docker wait test_container // 使用 docker wait 命令等待容器退出并获取退出代码

# 导出容器文件系统
docker export mycontainer // 将名为mycontainer的容器的文件系统导出到标准输出
docker export my_container > my_container_backup.tar // 将容器 my_container 的文件系统导出并保存到 my_container_backup.tar 文件中
docker export -o my_container_backup.tar my_container // 将容器 my_container 的文件系统导出并保存到 my_container_backup.tar 文件中
# ---------- 实列开始 ----------
# 启动容器
docker run -d --name my_container ubuntu bash -c "echo hello > /hello.txt && sleep 3600"
# 导出容器的文件系统
docker export my_container > my_container_backup.tar
# 查看导出的tar文件内容
tar -tf my_container_backup.tar
# 导入文件系统到新的容器
cat my_container_backup.tar | docker import - my_new_image
# ---------- 实列结束 ----------
# docker export 只导出容器的文件系统，不包括 Docker 镜像的层、元数据或运行时信息。
# 如果容器正在运行，导出的文件系统将是容器当前状态的快照。
# 导出的 tar 文件可能会很大，具体取决于容器的文件系统大小。
# docker export 命令是一个有用的工具，用于将容器的文件系统导出为 tar 归档文件。这对于备份、迁移和分析容器的文件系统非常有用。通过使用 --output 选项，用户可以将导出内容保存为指定文件，方便管理和使用。

# 查看容器的端口映射信息
docker port mycontainer
docker port mycontainer 80 // 查看特定端口的映射
docker port mycontainer 80/tcp // 查看特定端口和协议的映射

# 实时显示Docker容器的资源使用情况，包括CPU、内存、网络I/O和块I/O
docker stats
docker stats -a // 显示所有容器（包括未运行的容器）的资源使用情况
docker stats -a --no-trunc // 不截断输出

# 更新Docker容器的资源限制
docker update -m 2g mycontainer // 设置容器的内存限制为2GB
docker update --memory-swap 3g mycontainer // 设置容器内存和交换空间（swap）的总限制，如果设置为-1，表示不限制交换空间
docker update --cpu-shares 2048 mycontainer // 设置容器的 CPU 优先级，相对值。默认为 1024，较大的值表示较高的优先级。该选项不会直接限制容器的 CPU 使用量，而是控制 CPU 资源分配的优先级
docker update --cpus 2 mycontainer // 设置容器使用的 CPU 核心数。这个选项可以限制容器最多使用的 CPU 核心数
docker update --cpu-period 50000 my_container // 设置 CPU 周期时间。用于配合 --cpu-quota 限制容器的 CPU 使用时间。单位是微秒（默认值：100000 微秒 = 100ms）
docker update --cpu-quota 25000 my_container // 设置容器在每个周期内可以使用的最大 CPU 时间。单位是微秒。需要与 --cpu-period 配合使用
docker update --blkio-weight 800 my_container // 设置块 I/O 权重（范围：10 到 1000），表示容器对磁盘 I/O 操作的优先级。默认值为 500
docker update --pids-limit 200 my_container // 设置容器可以使用的最大进程数
docker update --restart always my_container // 设置容器的重启策略（no、on-failure、always、unless-stopped）
# ************************************************** 文件系统 **************************************************
# 将容器保存为新镜像
docker commit mycontainer my-new-image
docker commit mycontainer my-new-image:latest // 指定标签
docker commit -a "dongshan" -m "added new features" mycontainer my-new-image:latest // 添加作者信息和提交信息
docker commit --pause=false my_container my_new_image // 在不暂停容器的情况下，将其保存为新镜像
# ---------- 实列开始 ----------
# 启动一个容器
docker run -d -it --name my_container ubuntu bash
# 进行一些更改
docker exec my_container apt-get update
docker exec my_container apt-get install -y nginx
# 提交容器为新镜像
docker commit -a "Your Name" -m "Installed nginx" my_container my_new_image
# 查看新镜像
docker images
# ---------- 实列结束 ----------

# 从容器复制文件到宿主主机
docker cp my_container:/path/in/container/hello.txt /path/on/host/hello.txt

# 从容器复制目录到宿主主机
docker cp my_container:/path/in/container /path/on/host
docker cp -a my_container:/path/in/container /path/on/host // 保留所有元数据（权限、时间戳等）

# 从宿主主机复制文件到容器
docker cp /path/on/host/hello.txt my_container:/path/in/container/back.txt

# 从宿主主机复制目录到容器
docker cp /path/on/host my_container:/path/in/container
docker cp -a /path/on/host my_container:/path/in/container // 保留所有元数据（权限、时间戳等）

# 显示Docker容器文件系统的变更
docker diff my_container
# docker diff 命令的输出包含以下三种类型的变更：
    A: 表示新增的文件或目录
    D: 表示删除的文件或目录
    C: 表示修改过的文件或目录
# ************************************************** 版本信息 **************************************************
# 显示Docker系统的详细信息
docker info

# 显示Docker客户端和服务端的版本信息
docker version
docker --version // 只显示客户端版本的简要信息
# ************************************************** 网络命令 **************************************************
# 显示所有网络信息
docker network ls

# 查看指定网络的详细信息
docker network inspect network_name

# 创建一个新网络
docker network create my_network
# 常用参数
    --driver: 指定网络驱动程序（如 bridge、host、overlay）
    --subnet: 指定子网
    --gateway: 指定网关
    --ip-range: 指定可用 IP 地址范围
    --ipv6: 启用 IPv6
    --label: 为网络添加标签
# 示例
docker network create --driver bridge --subnet 192.168.1.0/24 my_network

# 删除一个或多个网络
docker network rm network1 network2

# 将一个容器连接到一个网络
docker network connect my_network my_container

# 将一个容器从一个网络断开
docker network disconnect my_network my_container
# ************************************************** 容器的卷 **************************************************
# 列出所有的卷
docker volume ls

# 查看某个卷的详细信息
docker volume inspect my_volume

# 创建一个新卷
docker volume create my_volume
# 常用参数
    --driver: 指定卷驱动程序（默认为 local）
    --label: 为卷添加标签
    -o, --opt: 为卷指定驱动程序选项
# 示例
docker volume create --name my_volume --label project=my_project

# 删除一个或多个卷
docker volume rm my_volume

# 删除未使用的卷
docker volume prune
# ************************************************** 容器编排 **************************************************
# docker compose run 命令用于启动一个新容器并允许一个特定的服务，而不启动整个Compose文件中定义的所有服务
# docker compose run 命令允许你在单个服务上执行任务，如运行一次性命令或调试
# 与docker compose up 的区别在于，run命令只会运行指定的服务，不会启动依赖它的其他服务

# 运行一个特定服务的容器
docker compose run web python manage.py migrate // 在名为web的容器中执行"python manage.py migrate"命令，而不启动其他服务

# 自动删除容器
docker compose run --rm web bash // 运行web容器，并启动一个Bash终端，任务完成后会自动删除容器

# 删除已停止的服务容器
docker compose rm

# 强制删除所有已停止的服务容器
docker compose rm -f

# 删除特定服务的已停止的容器
docker compose rm web

# 删除并清理相关卷
docker compose rm -v

# 列出所有运行中的容器
docker compose ps

# 列出所有容器（包括已停止的容器）
docker compose ps -a

# 仅显示容器ID
docker compose ps -q

# 列出所有服务的名称
docker compose ps --services

# 根据状态过滤容器
docker compose ps --filter "status=running"


# ************************************************** Docker完全卸载 **************************************************
1. 停止Docker服务
    sudo systemctl stop docker

2. 查看已安装的Docker软件包
    sudo yum list installed | grep docker

3. 卸载已安装的Docker软件包
    sudo yum remove containerd.io.x86_64 docker* -y

4. 删除Docker数据和配置文件
    sudo rm -rf /var/lib/docker         # 存放容器、镜像、卷、网络的配置
    sudo rm -rf /var/lib/container      # 管理Docker容器生命周期的组件（Docker容器的运行环境）
    sudo rm -rf /etc/docker             # Docker的配置文件

# ************************************************** Docker **************************************************




