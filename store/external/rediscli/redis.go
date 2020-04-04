package rediscli

import (
    "seqsvr/base/lib/logger"
)
import "github.com/go-redis/redis/v7"

var Cli *redis.Client

func Init() {
    Cli = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    pong, err := Cli.Ping().Result()
    logger.Infof("init redis %s, %v", pong, err)
}
