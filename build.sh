#!/bin/bash
# 当前目录
export BASE=$(pwd)
if [ ! -d "bin" ];then
mkdir bin
else
rm -rf bin/*
fi
echo "编译baas-gateway"
cd $BASE/baas-gateway
go build .
mv baas-gateway $BASE/bin
echo "编译baas-fabricengine"
cd $BASE/baas-fabricengine
go build .
mv baas-fabricengine $BASE/bin
echo "编译baas-kubeengine"
cd $BASE/baas-kubeengine
go build .
mv baas-kubeengine $BASE/bin
echo "编译baas-frontend"
cd $BASE/baas-frontend
rm -rf node_modules && npm install --registry=https://registry.npm.taobao.org
npm run build:prod
mv dist $BASE/bin/baas-frontend
