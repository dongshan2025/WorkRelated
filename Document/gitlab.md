
git pull gitlab/gitlab-ce:18.5.4-ce.0

docker run -d -p 8080:80 -p 8443:443 -p 8022:22 --restart always --name gitlab \
-v /usr/local/gitlab/config:/etc/gitlab \
-v /usr/local/gitlab/logs:/var/log/gitlab \
-v /usr/local/gitlab/data:/var/opt/gitlab \
--privileged=true gitlab/gitlab-ce:18.5.4-ce.0

在浏览器中输入：http://1.15.88.194:8080
等待初始化完成
完成后输入用户名：root 密码：<PASSWORD>

cd /usr/local/gitlab/config
vim gitlab.rb
在最上面添加以下内容：
external_url 'http://1.15.88.194'
gitlab_rails['gitlab_ssh_host'] = '1.15.88.194'
gitlab_rails['time_zone'] = 'Asia/Shanghai'

查看初始化密码：
cat initial_root_password

拉取代码(注意要把8080端口加上)：
git clone http://1.15.88.194:8080/dongshan/jinshu_backend.git

