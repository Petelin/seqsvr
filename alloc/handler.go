//go:generate zsh -c "protoc -I pb/ --go_out=plugins=grpc:./pb ./pb/*.proto"

package main

import (
	"context"
	"seqsvr/alloc/core"
	allocsvr "seqsvr/alloc/pb"
	"seqsvr/base/lib/grpcerr"

	"google.golang.org/grpc/codes"
)

type Server struct {
	service *core.Service
}

func (s *Server) FetchNextSeqNum(ctx context.Context, req *allocsvr.Uid) (*allocsvr.SeqNum, error) {
	seqNum, err := s.service.FetchNextSeqNum(req.GetUid(), req.GetVersion())
	resp := &allocsvr.SeqNum{SeqNum: seqNum}
	if err != nil {
		switch err {
		case core.ErrNotFoundUid:
			// client retry
			return nil, grpcerr.New(codes.NotFound, err.Error())
		case core.ErrMigrate:
			// rpc retry
			return nil, grpcerr.New(codes.Unavailable, err.Error())
		case core.ErrVersion:
			// success and retry
			resp.Version = s.service.RVersion
			resp.Router = make(map[string]*allocsvr.SectionIdArr)
			for k, v := range s.service.Rounter {
				tmp := &allocsvr.SectionIdArr{
					Val: v,
				}
				resp.Router[k] = tmp
			}
			return resp, nil
		}
		// rpc retry
		return nil, grpcerr.New(codes.Internal, err.Error())
	}
	return resp, nil
}
