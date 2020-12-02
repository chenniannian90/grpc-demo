#  双向证书认证


## 生成证书
### (服务器证书CN=server.io 会保错 "transport: authentication handshake failed: x509: certificate is valid for server.grpc.io, not server.io")
### 这里把CN 改为 server.grpc.io
```shell script
openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 3650 -subj "/C=GB/L=China/O=gobook/CN=github.com" -key ca.key -out ca.crt
	
openssl genrsa -out server.key 2048
openssl req -new -subj "/C=GB/L=China/O=server/CN=server.grpc.io" -key server.key -out server.csr
openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 -in server.csr -out server.crt
	
openssl genrsa -out client.key 2048
openssl req -new -subj "/C=GB/L=China/O=client/CN=client.grpc.io" -key client.key -out client.csr
openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.crt
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

## 证书拷贝
 ```shell script
cp ./server.crt ./server.key ./ca.crt  ./client.key ./client.crt  /etc/nginx/ssl/
```

## nginx 配置修改(grpcs: server 端处理加密数据)
```shell script
server {
   listen       443 ssl http2;
   ssl_certificate  /etc/nginx/ssl/server.crt;
   ssl_certificate_key /etc/nginx/ssl/server.key;
   ssl_client_certificate /etc/nginx/ssl/ca.crt;
   ssl_verify_client on;

   location / {
        grpc_ssl_certificate /etc/nginx/ssl/client.crt;
        grpc_ssl_certificate_key /etc/nginx/ssl/client.key;
        grpc_pass grpcs://localhost:1234;
   }
  }
```

## 重启 nginx
```shell script
nginx -s reload
```

## 修改 client.go line:35
```shell script
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

## nginx 配置修改(grpc: server 端处理明文数据)
```shell script
listen       443 ssl http2;
   ssl_certificate  /etc/nginx/ssl/server.crt;
   ssl_certificate_key /etc/nginx/ssl/server.key;
   ssl_client_certificate /etc/nginx/ssl/ca.crt;
   ssl_verify_client on;

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
