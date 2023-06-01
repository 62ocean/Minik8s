
# $1 函数文件所在目录
# $2 镜像名称
# $3 函数名称

# 生成所需的requirements.txt
pipreqs $1 --encodin utf8 --force

# build image
docker build -t $2 $1

# push image to docker hub
docker push $2

echo "--------- [$3] ------------"