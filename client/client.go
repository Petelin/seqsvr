package client

import (
	"context"
	"errors"
	"math/rand"
	allocsvr "seqsvr/alloc/pb"
	"seqsvr/base/common"
	"seqsvr/base/lib/logger"
	"seqsvr/base/lib/metricli"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

const retryPolicy = `{
  "methodConfig": [
    {
      "name": [
        {
          "service": "allocsvr.AllocServer"
        }
      ],
      "waitForReady": true,
      "retryPolicy": {
        "MaxAttempts": 10,
        "InitialBackoff": "0.002s",
        "MaxBackoff": "0.01s",
        "BackoffMultiplier": 1,
        "RetryableStatusCodes": [
          "UNAVAILABLE",
          "INTERNAL"
        ]
      }
    }
  ]
}`

var initAddr = []string{
	"127.0.0.1:9001",
	"127.0.0.1:9002",
}

type Client struct {
	sync.Mutex
	connPool map[string]allocsvr.AllocServerClient

	routeMap map[string][]uint64
	RVersion uint64
}

func NewClient(hosts []string) (*Client, error) {
	if len(hosts) == 0 {
		return nil, errors.New("host is empty")
	}
	i := rand.Intn(len(hosts))

	logger.Infof("random choice %s", hosts[i])

	conn, err := newConn(hosts[i])
	if err != nil {
		return nil, err
	}

	c := &Client{
		connPool: map[string]allocsvr.AllocServerClient{
			hosts[i]: conn,
		},
	}

	c.forceUpdateRouter()
	return c, nil
}

func (c *Client) forceUpdateRouter() {
	var cli allocsvr.AllocServerClient
	for _, v := range c.connPool {
		cli = v
		break
	}
	if cli == nil {
		return
	}

	resp, err := cli.FetchNextSeqNum(context.Background(), new(allocsvr.Uid))
	if err != nil {
		logger.Fatalf("rebuildRouter failed because FetchNextSeqNum failed, err=%v", err)
		return
	}

	c.updateRouterStatus(resp.GetVersion(), resp.GetRouter())
}

func (c *Client) updateRouterStatus(version uint64, router map[string]*allocsvr.SectionIdArr) {
	c.Lock()
	defer c.Unlock()
	if version <= c.RVersion {
		return
	}

	newPool := make(map[string]allocsvr.AllocServerClient, len(router))
	newRouterMap := make(map[string][]uint64)
	for name, ids := range router {
		newRouterMap[name] = ids.GetVal()
		if _, ok := c.connPool[name]; !ok {
			conn, err := newConn(name)
			if err != nil {
				continue
			}
			newPool[name] = conn
		} else {
			newPool[name] = c.connPool[name]
		}
	}
	c.routeMap = newRouterMap
	c.connPool = newPool
	c.RVersion = version
}

func newConn(add string) (allocsvr.AllocServerClient, error) {
	conn, err := grpc.Dial(add,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(retryPolicy),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  time.Millisecond,
				Multiplier: 1.6,
				Jitter:     0.1,
				MaxDelay:   time.Millisecond * 10,
			},
			MinConnectTimeout: time.Millisecond * 5,
		}))

	if err != nil {
		return nil, err
	}
	return allocsvr.NewAllocServerClient(conn), nil
}

func (c *Client) FetchNextSeqNum(ctx context.Context, entityID uint32) uint64 {
	start := time.Now()
	defer func() {
		d := time.Now().Sub(start)
		metricli.Histogram("alloc:FetchNextSeqNum", d.Microseconds())
	}()

	sID := common.GetSectionIDByUid(uint64(entityID))
	name := c.getServiceNameBySectionID(uint64(sID))
	if name == "" {
		// 路由表有问题
		c.forceUpdateRouter()
		return c.FetchNextSeqNum(ctx, entityID)
	}

	rpcCli, ok := c.connPool[name]
	if !ok {
		// 链接池子里没有, 不可能发生
		panic(ok)
		return c.fallBack(0)
	}

	resp, err := rpcCli.FetchNextSeqNum(ctx, &allocsvr.Uid{
		Uid:     uint64(entityID),
		Version: c.RVersion,
	})
	if err != nil {
		// err -> other error
		logger.Fatalf("rpc call failed, err=%v", err)
		return c.fallBack(0)
	}
	if resp == nil {
		return c.fallBack(0)
	}

	if resp.Router != nil {
		c.updateRouterStatus(resp.GetVersion(), resp.GetRouter())
		return c.FetchNextSeqNum(ctx, entityID)
	}
	if resp.GetSeqNum() == 0 {
		logger.Fatalf("0 , %v", resp)
	}
	return resp.GetSeqNum()
}

func (c *Client) getServiceNameBySectionID(sid uint64) string {
	for k, v := range c.routeMap {
		for _, id := range v {
			if id == sid {
				return k
			}
		}
	}

	return ""
}

func (c *Client) fallBack(i int) uint64 {
	return 0
}
