docker build -t ws:v1 -f Dockerfile .


docker run --name ws -p 9991:9991 -p 9992:9992 \
-v /home/dongshan/ws/etc/ws.yaml:/app/etc/ws.yaml \
-v /home/dongshan/ws/logs:/app/logs \
ws:v1



