//go:generate zsh -c "protoc -I pb/ --go_out=plugins=grpc:./pb ./pb/*.proto"

package main

import (
	"context"
	"seqsvr/alloc/core"
	allocsvr "seqsvr/alloc/pb"
	"seqsvr/base/lib/grpcerr"
	"seqsvr/base/lib/logger"

	"google.golang.org/grpc/codes"
)

type Server struct {
	service *core.Service
}

func (s *Server) FetchNextSeqNum(ctx context.Context, req *allocsvr.Uid) (*allocsvr.SeqNum, error) {
	logger.Infof("FetchNextSeqNum: %v", req)
	seqNum, b, err := s.service.FetchNextSeqNum(req.GetUid(), req.GetVersion())
	if err != nil {
		if err == core.ErrMigrate {
			return nil, grpcerr.New(codes.Unavailable, err.Error(), 14)
		}
	}
	resp := &allocsvr.SeqNum{SeqNum: seqNum}
	if b {
		resp.Version = s.service.RVersion
		resp.Router = make(map[string]*allocsvr.SectionIdArr)
		for k, v := range s.service.Rounter {
			tmp := &allocsvr.SectionIdArr{
				Val: v,
			}
			resp.Router[k] = tmp
		}
	}
	return resp, nil
}
