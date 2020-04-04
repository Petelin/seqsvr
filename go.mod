module seqsvr

go 1.13

replace protobuf/storesvr => ./protobuf/storesvr

require (
	github.com/go-redis/redis/v7 v7.2.0
	github.com/golang/protobuf v1.3.5
	go.uber.org/zap v1.14.1
	google.golang.org/grpc v1.28.0
)
