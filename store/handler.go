package main

import (
	"context"
	"fmt"
	"seqsvr/base/lib/logger"
	"seqsvr/store/external/rediscli"
	storesvr "seqsvr/store/pb"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v7"
)

const (
	MapRouterKey        = "maprouter"
	MapRouterVersionKey = "maprouter:version"
)

type RPCService struct {
}

func (R RPCService) UpdateMaxSeq(ctx context.Context, req *storesvr.UpdateMaxSeqReq) (*storesvr.NoContent, error) {
	logger.Infof("UpdateMaxSeq: %v", req)
	rediscli.Cli.Set(seqKey(req.GetSectionId()), req.GetMaxSeq(), -1)
	return &storesvr.NoContent{}, nil

}

func (R RPCService) GetSeqMax(ctx context.Context, req *storesvr.GetSeqMaxReq) (*storesvr.GetSeqMaxResp, error) {
	logger.Infof("GetSeqMax: %v", req)
	maxSeq, err := rediscli.Cli.Get(seqKey(req.GetSectionId())).Uint64()
	if err == redis.Nil {
		err = nil
	}
	return &storesvr.GetSeqMaxResp{
		MaxSeq: maxSeq,
	}, err
}

func (R RPCService) GetMapRouter(ctx context.Context, req *storesvr.GetMapRouterReq) (*storesvr.GetMapRouterResp, error) {
	logger.Infof("GetMapRouter: %v", req)

	var rawMap map[string]string
	var version uint64
	err := rediscli.Cli.Watch(func(tx *redis.Tx) error {
		var err error
		rawMap, err = rediscli.Cli.HGetAll(MapRouterKey).Result()
		if err != nil {
			return err
		}
		version, err = rediscli.Cli.Get(MapRouterVersionKey).Uint64()
		return err
	}, MapRouterVersionKey)
	if err != nil {
		return nil, err
	}

	var result = make(map[string]*storesvr.Sections, 16)
	for host, ids := range rawMap {
		ls := strings.Split(ids, ",")
		section := make([]uint64, 0, len(ls))
		for _, item := range ls {
			id, err := strconv.ParseUint(item, 0, 64)
			if err != nil {
				logger.Infof("parse routermap failed, %v", item)
				return nil, err
			}
			section = append(section, id)
		}
		result[host] = &storesvr.Sections{SectionIds: section}
	}
	return &storesvr.GetMapRouterResp{Version: version, RouterMap: result}, nil
}

func (R RPCService) SetHostRouter(ctx context.Context, req *storesvr.SetHostRouterReq) (*storesvr.GetMapRouterResp, error) {
	logger.Infof("SetHostRouter: %v", req)
	sb := strings.Builder{}
	for _, id := range req.GetSections().GetSectionIds() {
		sb.WriteString(fmt.Sprint(id))
		sb.WriteString(",")
	}
	value := sb.String()
	if len(value) > 0 {
		value = value[:len(value)-1]
	}

	var version uint64
	err := rediscli.Cli.Watch(func(tx *redis.Tx) error {
		var err error
		if value == "" {
			err = rediscli.Cli.HDel(MapRouterKey, req.GetHostName()).Err()
			if err != nil {
				return err
			}
		} else {
			err = rediscli.Cli.HSet(MapRouterKey, req.GetHostName(), value).Err()
			if err != nil {
				return err
			}
		}
		version, err = rediscli.Cli.Incr(MapRouterVersionKey).Uint64()
		return err
	}, MapRouterVersionKey)
	if err != nil {
		return nil, err
	}

	return R.GetMapRouter(ctx, &storesvr.GetMapRouterReq{})
}

func seqKey(key uint64) string {
	return fmt.Sprintf("seq:%d", key)
}
