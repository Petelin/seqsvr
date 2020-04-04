package metricli

import (
	"log"
	"os"
	"time"

	"github.com/rcrowley/go-metrics"
)

func Init() {
	go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
}

func Count(name string, n int64) {
	c := metrics.GetOrRegister(name, metrics.NewCounter())
	c.(metrics.Counter).Inc(n)
}
