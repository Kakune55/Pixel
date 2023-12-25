# Pixel
用Go语言实现的轻量级图床Web服务程序
数据库采用SQLite3

可实现基本的图片上传下载管理功能

## 现有接口
### 公共接口
- /upload 图片上传 
- /login 身份认证
- /img 图片调用
- /img/mini 图片缩略图调用
- /info 图片链接信息
### 需要认证的接口
- /img/del 图片删除
- /info/list 图片库
- /idlist 列出所有图片id

## 部分接口文档
### `/img` `/del` `/info`
#### 请求方式
Get  
#### 请求参数
id 图片id

### `/img/mini`
#### 请求方式
Get  
#### 请求参数
id 图片id
size 图片横向大小 单位像素 0-1024

## 使用的外部库
"github.com/disintegration/imaging" MIT LICENSE

# 部署
## 直接部署
- 确保本机拥有go环境
- 确保网络良好
~~~bash  
git clone https://github.com/Kakune55/Pixel.git #克隆存储库
cd ./Pixel #进入工作目录
go build main.go #编译
./main #运行
~~~
## 使用Docker (仅支持Linux-X86_64)
### 导入镜像
~~~bash
wget https://github.com/Kakune55/Pixel/releases/download/v1.0/pixel.tar.gz #版本号仅供参考
tar -zxvf pixel.tar.gz #解压
docker load pixel.tar #导入镜像
~~~
### 运行容器
~~~bash
docker run -d -p <你需要的端口号>:9090 --name:pixel pixel
~~~