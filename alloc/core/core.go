package core

import (
	"context"
	"errors"
	"seqsvr/base/common"
	"seqsvr/base/lib/logger"
	"seqsvr/base/lib/metricli"
	storesvr "seqsvr/store/pb"
	"sync"
	"sync/atomic"
	"time"
)

var ErrNotFoundUid = errors.New("not found the uid in the node")
var ErrVersion = errors.New("ErrVersion")
var ErrTimeout = errors.New("ErrTimeout")
var ErrMigrate = errors.New("ErrMigrate")
var ErrRouterNoChange = errors.New("ErrRouterNoChange")

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
	section map[common.SectionID]*common.Section

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
		section:     make(map[common.SectionID]*common.Section, 1000),
	}

	err := s.updateRouter()
	if err != nil {
		panic(err)
	}
	err = s.loadData(s.Rounter[s.name])

	go func() {
		t := time.NewTimer(time.Second * 1)
		for {
			select {
			case <-t.C:
				s.rMut.Lock()
				if s.updateRouter() != nil {
					s.rMut.Unlock()
					break
				}
				s.loadData(s.Rounter[s.name])
				s.rMut.Unlock()
			}

			t.Reset(time.Second)
		}
	}()
	return s
}

func (s *Service) updateRouter() error {
	resp, err := s.StoreClient.GetMapRouter(context.TODO(), &storesvr.GetMapRouterReq{})
	if err != nil {
		return err
	}

	if s.RVersion >= resp.GetVersion() {
		return ErrRouterNoChange
	}

	logger.Infof("update router (from %d,to %d) ...", s.RVersion, resp.GetVersion())

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
	metricli.Count("alloc:loadData", 1)
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

func (s *Service) FetchNextSeqNum(uid uint64, v uint64) (uint64, error) {
	s.rMut.RLock()
	defer s.rMut.RUnlock()
	metricli.Count("alloc:FetchNextSeqNum:req", 1)

	if s.stat == ServiceMigrate {
		return 0, ErrMigrate
	}

	// 拒绝老的请求，让他们去重试
	if s.RVersion > v {
		metricli.Count("alloc:FetchNextSeqNum:ErrVersion", 1)
		return 0, ErrVersion
	}

	var seqNum uint64
	sectionID := common.GetSectionIDByUid(uid)
	section, ok := s.section[sectionID]
	if !ok {
		return 0, ErrNotFoundUid
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
					return 0, nil
				}
				section.MaxSeq += common.Step
			}
			section.Mut.Unlock()
		} else {
			section.Mut.RUnlock()
		}
	}
	return seqNum, nil
}
