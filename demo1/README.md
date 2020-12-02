#  不加证书认证

## 编译
```shell script
go build client.go hello.pb.go
go build service.go hello.pb.go
````

# 运行
```shell script
./service
./client
```
 
 # nginx 配置修改
 ```shell script
server {
        listen       80 http2;
       
        location / {
           grpc_pass grpc://localhost:1234;

        }
    }
```
## 重启 nginx
```shell script
nginx -s reload
```
## 修改 client.go line:11
```go
conn, err := grpc.Dial("127.0.0.1:80", grpc.WithInsecure())
```
## 重新编译
```shell script
go build client.go hello.pb.go
```
## 再次运行(client 向nginx 发生请求, nginx 转发到 service)
```shell script
./service
./client
```