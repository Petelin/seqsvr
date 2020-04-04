//go:generate zsh -c "protoc -I pb/ --go_out=plugins=grpc:./pb ./pb/*.proto"

// Package main implements a server for Greeter service.
package main

import (
	"log"
	"net"
	"seqsvr/base/lib/logger"
	"seqsvr/store/external/rediscli"
	storesvr "seqsvr/store/pb"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8000"
)

func main() {
	logger.InitLogger(zap.NewDevelopmentConfig())
	rediscli.Init()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	storesvr.RegisterStoreServerServer(s, &RPCService{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
