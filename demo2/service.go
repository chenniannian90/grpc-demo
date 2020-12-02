package main

import (
	"context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

type HelloServiceImpl struct{}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *String,
) (*String, error) {
	reply := &String{Value: "hello:" + args.GetValue()}
	return reply, nil
}



func main()  {

	creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
