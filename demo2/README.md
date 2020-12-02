#  单向证书认证

## 生成证书
```shell script
openssl genrsa -out server.key 2048
openssl req -new -x509 -days 3650 -subj "/C=GB/L=China/O=grpc-server/CN=server.grpc.io" -key server.key -out server.crt
```


## 编译
```shell script
go build client.go hello.pb.go
go build service.go hello.pb.go
````

## 运行
```shell script
./service
./client
```

## 注意(如果是go版本在 1.15 及以上会报错,需要如下设置)
```shell script
export GODEBUG=x509ignoreCN=0
```
```shell script
rpc error: code = Unavailable desc = connection error: desc = "transport: authentication handshake failed: x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0"
参考 https://www.cnblogs.com/jackluo/p/13841286.html
```

 # 证书拷贝
 ```shell script
mkdir /etc/nginx/ssl/
cp ./server.crt /etc/nginx/ssl/
cp ./server.key /etc/nginx/ssl/
```
 
 ## nginx 配置修改(grpcs: server 端处理加密数据)
 ```shell script
server {
   listen       443 ssl http2;
   ssl_certificate  /etc/nginx/ssl/server.crt;
   ssl_certificate_key /etc/nginx/ssl/server.key;
   location / {
       grpc_pass grpcs://localhost:1234;
   }
  }
```
## 重启 nginx
```shell script
nginx -s reload
```
## 修改 client.go line:21
```go
conn, err := grpc.Dial("localhost:443", grpc.WithTransportCredentials(creds))
```
## 重新编译
```shell script
go build client.go hello.pb.go
```
## 再次运行(client 向nginx 发生请求, nginx 将加密数据转发到 service)
```shell script
./service
./client
```

 # nginx 配置修改(grpc: server 端处理明文数据)
 ```shell script
server {
   listen       443 ssl http2;
   ssl_certificate  /etc/nginx/ssl/server.crt;
   ssl_certificate_key /etc/nginx/ssl/server.key;
   location / {
       grpc_pass grpc://localhost:1234;
   }
  }
```
## 重启 nginx
```shell script
nginx -s reload
```
## 再次运行(client 向nginx 发生请求, nginx 将明文数据转发到 service)
```shell script
../demo1/service
./client
```
