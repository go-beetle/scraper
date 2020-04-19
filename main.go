package main

import (
	"container/heap"
	"github.com/go-beetle/scraper/scraper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func initZapLog() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()

	zap.ReplaceGlobals(logger)
}

func main() {
	initZapLog()
	// Some items and their priorities.
	now := int(time.Now().Unix())
	items := map[string]int{
		"google.com": now, "facebook.com": now+1, "reddit.com": now+2,
	}

	pq := make(scraper.PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &scraper.Url{
			Url:    value,
			ScrapeTime: priority,
			ScrapeIntervalSeconds: 5,
		}
		i += 1
	}
	heap.Init(&pq)
	for {
		zap.S().Debugf("Got URL %s", pq.GetNextUrl())
	}
}
