//go:generate zsh -c "protoc -I pb/ --go_out=plugins=grpc:./pb ./pb/*.proto"

package alloc

import (
	"context"
	allocsvr "seqsvr/alloc/pb"
	"seqsvr/base/lib/grpcerr"

	"google.golang.org/grpc/codes"
)

type Server struct {
	service *Service
}

func (s *Server) FetchNextSeqNum(ctx context.Context, req *allocsvr.Uid) (*allocsvr.SeqNum, error) {
	seqNum, b, err := s.service.FetchNextSeqNum(req.GetUid(), req.GetVersion())
	if err != nil {
		if err == ErrMigrate {
			return nil, grpcerr.New(codes.Unavailable, err.Error(), 14)
		}
	}
	rsp := &allocsvr.SeqNum{SeqNum: seqNum}
	if b {
		rsp.Version = s.service.RVersion
		rsp.Router = make(map[string]*allocsvr.SectionIdArr)
		for k, v := range s.service.Rounter {
			tmp := &allocsvr.SectionIdArr{
				Val: v,
			}
			rsp.Router[k] = tmp
		}
	}

	return rsp, nil
}
