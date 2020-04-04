package main

import (
	"context"
	"fmt"
	"seqsvr/alloc/core"
	"seqsvr/alloc/external"
	allocsvr "seqsvr/alloc/pb"
	"seqsvr/base/lib/logger"
	"testing"

	"go.uber.org/zap"
)

func TestHandler(t *testing.T) {
	logger.InitLogger(zap.NewDevelopmentConfig())
	ctx := context.TODO()
	s := &Server{
		service: core.NewService(ctx, "127.0.0.1:9001", external.GetStoreCli()),
	}
	fmt.Println(s.service.Rounter)

	for i := 0; i < 10; i++ {
		resp, err := s.FetchNextSeqNum(context.TODO(), &allocsvr.Uid{
			Uid:     1,
			Version: 2,
		})
		fmt.Println(resp.GetSeqNum(), err)
	}
}
