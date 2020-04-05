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

func Histogram(name string, n int64) {
	c := metrics.GetOrRegister(name, metrics.NewHistogram(metrics.NewUniformSample(5)))
	c.(metrics.Histogram).Update(n)
}
