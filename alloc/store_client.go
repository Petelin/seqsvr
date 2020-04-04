package alloc

import (
	storesvr "seqsvr/store/pb"

	"google.golang.org/grpc"
)

func GetStoreCli() storesvr.StoreServerClient {
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return storesvr.NewStoreServerClient(conn)
}
