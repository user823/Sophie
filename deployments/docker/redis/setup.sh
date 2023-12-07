sudo chmod 777 data
sudo chmod 777 log
sudo docker run -d -p 6379:6379 -v ./conf:/usr/local/etc/redis -v ./data:/data -v ./log:/var/log --name sophie-redis redis redis-server /usr/local/etc/redis/redis.conf