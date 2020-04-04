package client

import (
	"context"
	"seqsvr/base/common"
	"seqsvr/base/lib/logger"
	"seqsvr/base/lib/metricli"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

var cli *Client

func init() {
	logger.InitLogger(zap.NewDevelopmentConfig())
	var err error
	cli, err = NewClient([]string{"127.0.0.1:9001"})
	if err != nil {
		panic(err)
	}
}

// example
func TestOnce(t *testing.T) {
	println(cli.FetchNextSeqNum(context.Background(), 0))
	println(cli.FetchNextSeqNum(context.Background(), uint32(common.PerSectionIdSize)))
	println(cli.FetchNextSeqNum(context.Background(), uint32(common.PerSectionIdSize*2)))
}

func TestPresure(t *testing.T) {

	metricli.Init()
	wg := new(sync.WaitGroup)
	var uid = 0
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 10000000; j++ {
				ctx := context.Background()
				ctx, _ = context.WithTimeout(ctx, time.Millisecond*100)
				if cli.FetchNextSeqNum(ctx, uint32(i)) == 0 {
					metricli.Count("Err", 1)
				}
			}
		}(uid)
		uid += int(common.PerSectionIdSize)
	}
	wg.Wait()
}
