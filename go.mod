module seqsvr

go 1.13

replace protobuf/storesvr => ./protobuf/storesvr

require (
	github.com/cyberdelia/go-metrics-graphite v0.0.0-20161219230853-39f87cc3b432
	github.com/go-redis/redis/v7 v7.2.0
	github.com/golang/protobuf v1.3.5
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0
	go.uber.org/atomic v1.6.0
	go.uber.org/zap v1.14.1
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	google.golang.org/grpc v1.28.0
)
