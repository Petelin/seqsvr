package client

import (
	"context"
	"fmt"
	"math/rand"
	"seqsvr/base/common"
	"seqsvr/base/lib/logger"
	"seqsvr/base/lib/metricli"
	"sync"
	"testing"
	"time"

	"go.uber.org/atomic"
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
			var beforeVal uint64 = 0
			defer wg.Done()
			for j := 0; j < 10000000; j++ {
				ctx := context.Background()
				ctx, _ = context.WithTimeout(ctx, time.Millisecond*100)
				newVal := cli.FetchNextSeqNum(ctx, uint32(i))
				if newVal == 0 {
					metricli.Count("client:Err:failed", 1)
				} else if newVal <= beforeVal {
					metricli.Count("client:Err:Sequence", 1)
				}
			}
		}(uid)
		uid += int(common.PerSectionIdSize)
	}
	wg.Wait()
}

func TestBenchmark(t *testing.T) {
	var c atomic.Uint64

	wg := new(sync.WaitGroup)
	startT := time.Now()
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100000; i++ {
				cli.FetchNextSeqNum(context.TODO(), rand.Uint32()%3)
				c.Add(1)
			}
		}()
	}
	wg.Wait()
	d := time.Now().Sub(startT)
	fmt.Println(d, c.Load(), c.Load()/uint64(d/time.Second))
	// result: 16线程 24.011905654s 1600000 66666
	// 		    8线程 19.526373785s  800000 42105
}

func TestRateLimit(t *testing.T) {
	metricli.Init()

	rate := time.Second / 80000
	burstLimit := 100
	tick := time.NewTicker(rate)
	defer tick.Stop()
	throttle := make(chan time.Time, burstLimit)
	go func() {
		for t := range tick.C {
			select {
			case throttle <- t:
			default:
			}
		} // does not exit after tick.Stop()
	}()
	wg := new(sync.WaitGroup)
	for range make([]int, 100000000) {
		<-throttle // rate limit our Service.Method RPCs
		wg.Add(1)
		go func() {
			cli.FetchNextSeqNum(context.TODO(), 0)
			wg.Done()
		}()
	}
	wg.Wait()
}
