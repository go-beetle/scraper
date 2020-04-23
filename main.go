package main

import (
	"bytes"
	"github.com/go-beetle/scraper/scraper"
	"github.com/go-beetle/scraper/scraper/cache"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initZapLog() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level.SetLevel(zap.InfoLevel)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()

	zap.ReplaceGlobals(logger)
}

func main() {
	initZapLog()
	// Some items and their priorities.
	urls := []*scraper.ScraperItem{
		{
			Url:                       &scraper.URL{"worldometers.info", "/"},
			ScrapeIntervalNanoseconds: 5 * 1e9,
		},
	}

	pq := scraper.InitPQSeed(urls)
	htmlCache := cache.InitHTMLCache(10 * 1024 * 1024)

	for {
		scraperItem := pq.PeekScraperItemAndUpdate()
		url := scraperItem.Url

		body, err := scraper.Get(url.String())
		if err != nil {
			zap.S().Error("Error with scraper get url")
		}

		hrefs := scraper.GetHref(string(body), url)
		zap.S().Debugf("Processed found %d links with hostname %s", len(hrefs), url.Hostname)

		pq.AddURLs(hrefs, scraperItem)

		prevBody := htmlCache.Get(url.Concat())
		if prevBody == nil {
			htmlCache.Add(url.Concat(), body, int(scraperItem.ScrapeIntervalNanoseconds))
			scraper.WriteFile(body, url)

		} else if bytes.Compare(body, prevBody) != 0 {
			scraper.WriteFile(body, url)
		}

	}
}
