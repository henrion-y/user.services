#!/bin/bash

# 镜像仓库地址
REGISTRY=""
NAMESPACE=""
SERVICE_NAME="user.services"

# 获取当前版本号
VERSION=$(git describe --tags --always)

# 构建镜像
docker build -t "$REGISTRY/$NAMESPACE/$SERVICE_NAME:$VERSION" .
docker tag "$REGISTRY/$NAMESPACE/$SERVICE_NAME:$VERSION" "$REGISTRY/$NAMESPACE/$SERVICE_NAME:latest"

# 推送镜像
docker push "$REGISTRY/$NAMESPACE/$SERVICE_NAME:$VERSION"
docker push "$REGISTRY/$NAMESPACE/$SERVICE_NAME:latest"
