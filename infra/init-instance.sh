#!/bin/bash
sudo yum update -y
sudo yum install docker -y
sudo systemctl start docker
sudo docker pull mysql:latest
sudo usermod -a -G docker ec2-user
id ec2-user
newgrp docker
docker run --name my-mysql -e MYSQL_ROOT_PASSWORD=my-secret-password -p 3306:3306 -d mysql:latest
docker run -d --name redis-server -p 6379:6379 redis
