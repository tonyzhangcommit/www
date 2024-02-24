#!/bin/bash

# RabbitMQ Docker 容器启动脚本

# 创建一个数据卷，用于持久化数据
docker volume create rabbitmq_videosystem_data

# 启动 RabbitMQ 容器
docker run -d \
    --name rabbitmq \
    -p 5672:5672 \
    -p 15672:15672 \
    -v rabbitmq_videosystem_data:/var/lib/rabbitmq \
    --hostname vsmrabbitmqhost \
    --restart always \
    rabbitmq:latest