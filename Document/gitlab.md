
git pull gitlab/gitlab-ce:18.5.4-ce.0

docker run -d -p 9080:80 -p 9443:443 -p 9022:22 --restart always --name gitlab \
-v /usr/local/gitlab/config:/etc/gitlab \
-v /usr/local/gitlab/logs:/var/log/gitlab \
-v /usr/local/gitlab/data:/var/opt/gitlab \
--privileged=true gitlab/gitlab-ce:18.5.4-ce.0

Docker启动之后：
cd /usr/local/gitlab/config
vim gitlab.rb
在最上面添加以下内容：
external_url 'http://1.15.88.194'

然后使用如下命令进入容器：
docker exec -it gitlab /bin/bash
cd cd /opt/gitlab/embedded/service/gitlab-rails/config
vi gitlab.yml
将端口80改为9080
保存退出

在容器内执行：gitlab-ctl restart
重启后退出容器

进入/usr/local/gitlab/config目录
cd /usr/local/gitlab/config
查看初始化密码：
cat initial_root_password

在浏览器中输入：http://1.15.88.194:9080
等待初始化完成
完成后输入用户名：root 密码：<Dongshan123456>

拉取代码(注意要把8080端口加上)：
git clone http://1.15.88.194:9080/dongshan/jinshu_backend.git

