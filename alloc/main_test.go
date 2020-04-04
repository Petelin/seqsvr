package main

import (
	"context"
	"fmt"
	"seqsvr/alloc/core"
	"seqsvr/alloc/external"
	allocsvr "seqsvr/alloc/pb"
	"testing"
)

func TestHandler(t *testing.T) {
	ctx := context.TODO()
	s := &Server{
		service: core.NewService(ctx, "127.0.0.1:9001", external.GetStoreCli()),
	}
	fmt.Println(s.service.Rounter)

	for i := 0; i < 1000; i++ {
		resp, err := s.FetchNextSeqNum(context.TODO(), &allocsvr.Uid{
			Uid:     1,
			Version: 2,
		})
		fmt.Println(resp.GetSeqNum(), err)
	}
}
