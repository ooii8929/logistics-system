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
ECR_REPOSITORY=${ECR_REGISTRY}/logistics-track
ECR_LOGIN_PASSWORD=$(aws ecr get-login-password --region ${AWS_REGION})
echo ${ECR_LOGIN_PASSWORD} | sudo docker login -u AWS --password-stdin ${ECR_REGISTRY}

sudo docker login -u AWS -p ${ECR_LOGIN_PASSWORD} ${ECR_REGISTRY}

# 拉取并运行MySQL和Redis
sudo docker pull mysql:latest
sudo docker pull redis:latest
docker run --name my-mysql -e MYSQL_ROOT_PASSWORD=my-secret-password -p 3306:3306 -d mysql:latest
docker run -d --name redis-server -p 6379:6379 redis

# 拉取并运行您的应用程序镜像
sudo docker pull ${ECR_REPOSITORY}:latest
docker run -d --name your-app-container \
           -p 8080:8080 \ # 调整端口映射，根据您的应用程序需要
           --link my-mysql:mysql \
           --link redis-server:redis \
           ${ECR_REPOSITORY}:latest

# 配置 Nginx
sudo bash -c 'cat > /etc/nginx/sites-available/default << EOL
server {
    listen 80;
    server_name _; # 修改为您的域名或IP

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
