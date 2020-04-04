package alloc

import (
	"context"
	"fmt"
	allocsvr "seqsvr/alloc/pb"
	"testing"
)

func TestHandler(t *testing.T) {
	ctx := context.TODO()
	s := &Server{
		service: NewService(ctx, "127.0.0.1:9001", GetStoreCli()),
	}
	fmt.Println(s.service.Rounter)

	for i := 0; i < 1000; i++ {
		fmt.Println(s.FetchNextSeqNum(context.TODO(), &allocsvr.Uid{
			Uid:     1,
			Version: 2,
		}))
	}
}
