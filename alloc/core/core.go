package core

import (
	"context"
	"errors"
	"seqsvr/base/common"
	storesvr "seqsvr/store/pb"
	"sync"
	"sync/atomic"
)

var ErrNotFoundUid = errors.New("not found the uid in the node")
var ErrVersion = errors.New("not found the uid in the node")
var ErrTimeout = errors.New("not found the uid in the node")
var ErrMigrate = errors.New("not found the uid in the node")

type Stat uint64

const (
	OnService Stat = iota
	PreMigrate
	Migrate
	Rollback
)

const (
	_ Stat = iota
	ServiceNormal
	ServiceMigrate
)

type Migration struct {
	Stat      Stat                          `json:"stat"`
	ChangedId []common.SectionID            `json:"changed_id"`
	Change    map[string][]common.SectionID `json:"change"`
	Version   uint64                        `json:"version"`
}

type Service struct {
	ctx context.Context

	name    string
	rMut    *sync.RWMutex
	section map[common.SectionID]common.Section

	nextRVersion uint64
	nextSection  map[common.SectionID]common.Section

	RVersion uint64
	Rounter  map[string][]uint64

	stat Stat

	StoreClient storesvr.StoreServerClient
}

func NewService(ctx context.Context, name string, client storesvr.StoreServerClient) *Service {
	s := &Service{
		ctx:         ctx,
		name:        name,
		rMut:        new(sync.RWMutex),
		stat:        OnService,
		StoreClient: client,
		section:     make(map[common.SectionID]common.Section, 1000),
	}

	err := s.updateRouter()
	if err != nil {
		panic(err)
	}
	err = s.loadData(s.Rounter[s.name])
	return s
}

func (s *Service) updateRouter() error {
	s.rMut.RLock()
	defer s.rMut.RUnlock()

	resp, err := s.StoreClient.GetMapRouter(context.TODO(), &storesvr.GetMapRouterReq{})
	if err != nil {
		return err
	}

	result := make(map[string][]uint64, len(resp.GetRouterMap()))
	if resp != nil {
		for k, v := range resp.GetRouterMap() {
			result[k] = v.GetSectionIds()
		}
	}
	s.RVersion = resp.GetVersion()
	s.Rounter = result
	return nil
}

func (s *Service) loadData(ids []uint64) error {
	s.rMut.RLock()
	defer s.rMut.RUnlock()

	for _, id := range ids {
		resp, err := s.StoreClient.GetSeqMax(context.TODO(), &storesvr.GetSeqMaxReq{
			SectionId: id,
		})
		if err != nil {
			return err
		}
		sid := common.SectionID(id)
		s.section[sid] = common.NewSection(sid, resp.GetMaxSeq())
	}
	return nil
}

func (s *Service) FetchNextSeqNum(uid uint64, v uint64) (uint64, bool, error) {
	s.rMut.RLock()
	defer s.rMut.RUnlock()

	if s.stat == ServiceMigrate {
		return 0, false, ErrMigrate
	}

	var routerChange bool
	if s.RVersion > v {
		routerChange = true
	}

	var seqNum uint64
	sectionID := common.GetSectionIDByUid(uid)
	section, ok := s.section[sectionID]
	if !ok {
		return 0, routerChange, ErrNotFoundUid
	}
	if ok {
		_, index := common.CalcIndex(section.RangeID, uid)
		section.Mut.RLock()
		seqNum = atomic.AddUint64(&section.SeqNum[index], 1)
		if seqNum > section.MaxSeq {
			section.Mut.RUnlock()
			section.Mut.Lock()
			if seqNum > section.MaxSeq {
				_, err := s.StoreClient.UpdateMaxSeq(context.TODO(), &storesvr.UpdateMaxSeqReq{
					SectionId: uint64(sectionID),
					MaxSeq:    section.MaxSeq + common.Step,
				})
				if err != nil {
					return 0, routerChange, nil
				}
				section.MaxSeq += common.Step
			}
			section.Mut.Unlock()
		} else {
			section.Mut.RUnlock()
		}
	}
	return seqNum, routerChange, nil
}
