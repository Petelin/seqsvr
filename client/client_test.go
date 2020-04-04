package client

import (
	"context"
	"seqsvr/base/lib/logger"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// example
func TestExample(t *testing.T) {
	logger.InitLogger(zap.NewDevelopmentConfig())
	cli, err := NewClient([]string{"127.0.0.1:9001"})
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, time.Millisecond*5)
		println(i, "--", cli.FetchNextSeqNum(ctx, uint32(i)))
	}
}
