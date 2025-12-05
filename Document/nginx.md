# Docker安装Nginx
docker run --name nginx -d -p 80:80 nginx:1.29.3-perl
// ro是 read-only 的缩写，表示将宿主机上的某个目录或卷以只读权限挂载到容器的指定目录。这允许容器读取文件，但不能修改它们
docker run --name nginx -d --restart always  -p 80:80 \
-v /home/dongshan/nginx/config/nginx.conf:/etc/nginx/nginx.conf:ro \
-v /home/dongshan/nginx/logs:/var/log/nginx \
nginx:1.29.3-perl

# 进入容器内
docker exec -it nginx bash

# 查看默认nginx.conf文件
cat /etc/nginx/nginx.conf

# 查看默认html文件
cat /usr/share/nginx/html/index.html

# 日志路径
/var/log/nginx

# ================================================== 反向代理配置 ==================================================
user  nginx;
worker_processes  auto;

error_log  /var/log/nginx/error.log notice;
pid        /run/nginx.pid;


events {
    worker_connections  1024;
}


http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;
    #tcp_nopush     on;

    keepalive_timeout  65;

    #gzip  on;

    include /etc/nginx/conf.d/*.conf;

    server {
        listen 9991;
        server_name 192.168.252.101;

        location /v1/ws/ {
            proxy_set_header   X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header   Host      $http_host;
            proxy_set_header X-NginX-Proxy true;
            proxy_pass http://192.168.252.103:9991/v1/ws;
            proxy_http_version         1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "Upgrade";
            proxy_redirect off;
        }
    }
}
# ================================================== WebSocket负载均衡 ==================================================
user  nginx;
worker_processes  auto;

error_log  /var/log/nginx/error.log notice;
pid        /run/nginx.pid;


events {
    worker_connections  1024;
}


http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;
    #tcp_nopush     on;

    keepalive_timeout  65;

    #gzip  on;

    include /etc/nginx/conf.d/*.conf;

    upstream ws_backend {
	    server 192.168.252.102:9991;
        server 192.168.252.103:9991;
    }

    server {
        listen 9991;
        server_name 192.168.252.101;

        location /v1/ws/ {
            proxy_set_header   X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header   Host      $http_host;
            proxy_set_header X-NginX-Proxy true;
            proxy_pass http://ws_backend/v1/ws;
            proxy_http_version         1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "Upgrade";
            proxy_redirect off;
        }
    }
}
# ==================================================  ==================================================

docker run --name nginx -p 80:80 -p 9991:9991 -v /home/dongshan/nginx/config/nginx.conf:/etc/nginx/nginx.conf nginx:1.29.3-perl





