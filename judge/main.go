package judge

import (
	"fmt"
	"seqsvr/base/lib/logger"
	"seqsvr/judge/external/rediscli"
	"strconv"
)

func init() {
	logger.InitLogger()
	rediscli.Init()
}

// alloc watch mapchanged
func WatchRouterChange(callback func(version uint64)) {
	pubsub := rediscli.Cli.Subscribe("rmchannel")
	c := pubsub.Channel()
	for {
		select {
		case message := <-c:
			x, _ := strconv.ParseUint(message.Payload, 10, 64)
			callback(x)
		}
	}
}

func NotifyRouterChange(version uint64) {
	fmt.Println("pusb")
	rediscli.Cli.Publish("rmchannel", fmt.Sprint(version))
}
