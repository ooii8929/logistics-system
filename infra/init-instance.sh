#!/bin/bash
sudo apt update -y
sudo apt install docker.io -y
sudo apt install awscli -y
sudo apt install nginx -y
sudo systemctl start docker
sudo usermod -a -G docker ubuntu
id ubuntu
newgrp docker

# Define AWS region
AWS_REGION="ap-northeast-1"

# Configure AWS CLI
aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin 730461323800.dkr.ecr.${AWS_REGION}.amazonaws.com

# Log in to ECR
ECR_REGISTRY=730461323800.dkr.ecr.${AWS_REGION}.amazonaws.com
ECR_REPOSITORY_APPLICATION=${ECR_REGISTRY}/logistics-track
ECR_REPOSITORY_DATABASE=${ECR_REGISTRY}/logistics-database
ECR_LOGIN_PASSWORD=$(aws ecr get-login-password --region ${AWS_REGION})
echo ${ECR_LOGIN_PASSWORD} | sudo docker login -u AWS --password-stdin ${ECR_REGISTRY}

sudo docker login -u AWS -p ${ECR_LOGIN_PASSWORD} ${ECR_REGISTRY}
docker network create my_app_network

# 拉取并运行MySQL和Redis

sudo docker pull redis:latest
sudo docker pull ${ECR_REPOSITORY_DATABASE}:latest
docker run -d --name logistics-db \
           --network my_app_network \
           -p 3306:3306 \
           ${ECR_REPOSITORY_DATABASE}:latest

pwd

echo "requirepass logisticsredis" > redis.conf
# docker run -d --name redis-server --network my_app_network -p 6379:6379 -v "$(pwd)/redis.conf:/usr/local/etc/redis/redis.conf" -d redis:latest redis-server /usr/local/etc/redis/redis.conf
docker run -d --name redis-server --network my_app_network -p 6379:6379 -v "$(pwd)/redis.conf:/usr/local/etc/redis/redis.conf" -d redis:latest redis-server /usr/local/etc/redis/redis.conf


# 拉取并运行您的应用程序镜像
sudo docker pull ${ECR_REPOSITORY_APPLICATION}:latest
docker run -d --name logistics-system \
           --network my_app_network \
           -p 8080:8080 \
           ${ECR_REPOSITORY_APPLICATION}:latest


# 配置 Nginx
sudo bash -c 'cat > /etc/nginx/sites-available/default << EOL
server {
    listen 80;
    server_name logistics.wooah.app;

    keepalive_timeout 600s;
    client_header_timeout 600s;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOL'

# 启动 Nginx
sudo systemctl enable nginx
sudo systemctl restart nginx

# For cronjob worker
# PHY_INTERFACE=$(ip -4 route list 0/0 | awk '{ print $5 }')
# IP_ADDRESS=$(ip -4 -o addr show dev $PHY_INTERFACE primary | awk '{gsub(/\/(.*)/, "", $4); print $4}')
# sudo docker network create -d macvlan \
#   --subnet=$SUBNET \
#   --gateway=$GATEWAY \
#   -o parent=$PHY_INTERFACE \
#   mymacvlan

# sudo docker network connect mymacvlan my-mysql

# 创建 call_generate_report.sh 文件
# 创建 call_generate_report.sh 文件
cat > call_generate_report.sh << EOL
#!/bin/bash

curl -X GET https://logistics.wooah.app/generate_report
EOL

# 将脚本移动到 /usr/local/bin
sudo mv call_generate_report.sh /usr/local/bin

# 为 call_generate_report.sh 添加可执行权限
sudo chmod +x /usr/local/bin/call_generate_report.sh

# 添加 Cronjob
(crontab -l ; echo "0 0,8,16 * * * /usr/local/bin/call_generate_report.sh") | crontab -
