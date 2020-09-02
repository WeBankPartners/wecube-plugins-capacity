# Wecube-Plugins-Capacity安装指引

## docker运行：
  
#### 如果已存在docker镜像，请忽略1、2、4步
#### 1. 下载release上最新的wecube-plugins-capacity.zip压缩包，装好所依赖docker和mysql
#### 2. 解压wecube-plugins-capacity.zip压缩包，里面会有image.tar镜像和init.sql
#### 3. 在mysql上导入init.sql初始化数据库，如果是zip包解压的会有init.sql，如果是项目编译的，sql文件在 doc/init.sql 
#### 4. 导入docker镜像
```
docker load --input image.tar
```
#### 5. 创建本地目录，例如/app/test(如果是其它目录请替换下面命令中的/app/test)，替换如下命令的mysql连接参数
```
必填
CAPACITY_MYSQL_HOST=127.0.0.1 -> 把127.0.0.1替换成mysql地址
CAPACITY_MYSQL_PORT=3306  -> 把3306替换成mysql端口
CAPACITY_MYSQL_USER=root -> 把root替换成mysql用户
CAPACITY_MYSQL_PWD=wecube -> 把wecube替换成mysql用户密码
wecube-plugins-capacity:v0.0.1 -> 把后面的版本号改成所导入镜像的版本号
可选
GATEWAY_URL=http://127.0.0.1:8080 -> Open-Monitor的服务地址
```
```
mkdir -p /app/test
docker run --name capacity --volume /app/test/logs:/app/capacity/logs --volume /app/test/public/r_images:/app/capacity/public/r_images --volume /etc/localtime:/etc/localtime -d -p 9096:9096 -e CAPACITY_MYSQL_HOST=127.0.0.1 -e CAPACITY_MYSQL_PORT=3306 -e CAPACITY_MYSQL_USER=root -e CAPACITY_MYSQL_PWD=wecube -e CAPACITY_LOG_LEVEL=debug  wecube-plugins-capacity:v0.0.1
```
容器运行起来后打开 http://127.0.0.1:9096/capacity/ 界面