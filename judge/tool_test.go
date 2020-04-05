package judge

import (
	"context"
	"fmt"
	"seqsvr/alloc/external"
	"seqsvr/base/lib/logger"
	storesvr "seqsvr/store/pb"
	"testing"
)

func init() {
	logger.InitLogger()
}

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

func TestWatchRouter(t *testing.T) {
	WatchRouterChange(func(version uint64) {
		fmt.Println("change", version)
	})
}

func TestRouterNotify(t *testing.T) {
	NotifyRouterChange(2)
}
