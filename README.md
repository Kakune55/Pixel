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
### `/img` `/img/mini` `/del` `/info`
#### 请求方式
Get  
#### 请求参数
id 图片id

## 使用的外部库
"github.com/disintegration/imaging" MIT LICENSE
