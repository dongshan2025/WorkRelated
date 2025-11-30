# node-exporter安装
在192.168.252.102服务器上运行以下容器：

docker run -d \
  --net="host" \
  --pid="host" \
  -v "/:/host:ro,rslave" \
  prom/node-exporter:v1.10.2 \
  --path.rootfs=/host

启动容器后，在浏览器中查看当前服务器信息：
http://192.168.252.102:9100/
http://192.168.252.102:9100/metrics

# prometheus安装
在192.168.252.101服务器上运行以下容器：

启动一个临时容器，主要复制配置文件：
docker run --name prometheus -d -p 9090:9090 prom/prometheus:v3.7.3
docker cp prometheus:/etc/prometheus/prometheus.yml /prometheus

修改配置文件如下：
====================================================================================================
# my global config
global:
  scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "prometheus"

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
      - targets: ["localhost:9090"]
       # The label name is added as a label `label_name=<label_value>` to any timeseries scraped from this config.
        labels:
          app: "prometheus"
  - job_name: 'node'
    static_configs:
      - targets: ['192.168.252.102:9100']
====================================================================================================
添加的内容是以下内容：
  - job_name: 'node'
    static_configs:
      - targets: ['192.168.252.102:9100']

删除临时容器后，重新启动容器：
docker run --name prometheus -d -p 9090:9090 -v /home/dongshan/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus:v3.7.3

容器启动后，在浏览器中输入如下地址查看配置的信息：
http://192.168.252.101:9090/targets


# grafana安装
在192.168.252.101服务器上运行以下容器：

docker run -d --name=grafana -p 3000:3000 grafana/grafana:12.4.0-19363970803

在浏览器中输入以下地址登录：
http://192.168.252.101:3000/login
admin/admin
admin/NaDW6xNjBTaG4Gr

点击右侧"Dashboards"-"Import a dashboard" 然后点击以下的链接：
Find and import dashboards for common applications at grafana.com/dashboards

https://grafana.com/grafana/dashboards/
选择一个合适的dashboards，然后下载它的json文件，将该文件复制到"http://192.168.252.101:3000/dashboard/import"中的文本框中，然后点击"Load"按钮进行导入即可。

http://192.168.252.101:3000/d/rYdddlPWk/node-exporter-full?orgId=1&from=now-30m&to=now&timezone=browser&var-ds_prometheus=ff4hjoxmjfke8c&var-job=node&var-nodename=ubuntu&var-node=192.168.252.102:9100&refresh=10s



