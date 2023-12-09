#!/bin/bash
sudo yum update -y
sudo yum install docker -y
sudo yum install -y aws-cli
sudo systemctl start docker
sudo usermod -a -G docker ec2-user
id ec2-user
newgrp docker

# 配置 AWS CLI
export AWS_REGION=ap-northeast-1
aws configure set aws_access_key_id YOUR_AWS_ACCESS_KEY
aws configure set aws_secret_access_key YOUR_AWS_SECRET_KEY
aws configure set default.region ${AWS_REGION}

# 登录 ECR
ECR_REGISTRY=730461323800.dkr.ecr.ap-northeast-1.amazonaws.com
ECR_REPOSITORY=${ECR_REGISTRY}/logistics-track
ECR_LOGIN_PASSWORD=$(aws ecr get-login-password --region ${AWS_REGION})
sudo docker login -u AWS -p ${ECR_LOGIN_PASSWORD} ${ECR_REGISTRY}

# 拉取和运行MySQL和Redis
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
