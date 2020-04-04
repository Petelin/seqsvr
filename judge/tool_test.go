package judge

import (
	"context"
	"seqsvr/alloc/external"
	storesvr "seqsvr/store/pb"
	"testing"
)

func TestTransfer(t *testing.T) {
	cli := external.GetStoreCli()

	ctx := context.Background()
	cli.SetHostRouter(ctx, &storesvr.SetHostRouterReq{
		HostName: "127.0.0.1:9001",
		Sections: &storesvr.Sections{
			SectionIds: []uint64{0, 1},
		},
	})
	cli.SetHostRouter(ctx, &storesvr.SetHostRouterReq{
		HostName: "127.0.0.1:9002",
		Sections: &storesvr.Sections{
			SectionIds: []uint64{2},
		},
	})
}

func TestTransferB(t *testing.T) {
	cli := external.GetStoreCli()

	ctx := context.Background()
	cli.SetHostRouter(ctx, &storesvr.SetHostRouterReq{
		HostName: "127.0.0.1:9001",
		Sections: &storesvr.Sections{
			SectionIds: []uint64{0, 2},
		},
	})
	cli.SetHostRouter(ctx, &storesvr.SetHostRouterReq{
		HostName: "127.0.0.1:9002",
		Sections: &storesvr.Sections{
			SectionIds: []uint64{1},
		},
	})
}
