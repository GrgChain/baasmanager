#!/bin/bash
# 当前目录
export BASE=$(pwd)
if [ ! -d "bin" ];then
mkdir bin
else
echo "文件夹已经存在"
fi
# 编译baas-gateway
cd $BASE/baas-gateway
go build .
mv baas-gateway $BASE/bin
# 编译baas-fabricengine
cd $BASE/baas-fabricengine
go build .
mv baas-fabricengine $BASE/bin
# 编译baas-kubeengine
cd $BASE/baas-kubeengine
go build .
mv baas-kubeengine $BASE/bin