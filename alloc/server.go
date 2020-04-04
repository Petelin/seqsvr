package main

import (
	"context"
	"flag"
	"log"
	"net"
	"seqsvr/alloc/core"
	"seqsvr/alloc/external"
	allocsvr "seqsvr/alloc/pb"
	"seqsvr/base/lib/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := flag.String("port", "9001", "rpc port")
	flag.Parse()

	logger.InitLogger(zap.NewDevelopmentConfig())
	address := "127.0.0.1:" + *port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	allocsvr.RegisterAllocServerServer(s, &Server{
		service: core.NewService(context.Background(), address, external.GetStoreCli()),
	})

	logger.Infof("alloc start service: %s", address)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
