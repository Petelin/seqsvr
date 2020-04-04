package alloc

import (
    storesvr "seqsvr/store/pb"

    "google.golang.org/grpc"
)

var storeCli storesvr.StoreServerClient

func InitStoreCli()  {
    conn, err := grpc.Dial("localhost:8000")
    if err != nil{
        panic(err)
    }
    storeCli = storesvr.NewStoreServerClient(conn)
}
