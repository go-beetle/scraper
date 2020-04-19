package main

import (
	"github.com/go-beetle/scraper/scraper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	urls := []*scraper.ScraperItem{
		{
			Url:                       &scraper.URL{"worldometers.info", "/"},
			ScrapeIntervalNanoseconds: 5 * 1e9,
		},
	}

	pq := scraper.InitPQSeed(urls)

	for {
		scraperItem := pq.PeekScraperItemAndUpdate()
		url := scraperItem.Url

		body, err := scraper.Get(url.String())
		if err != nil {
			zap.S().Error("Error with scraper get url")
		}

		hrefs := scraper.GetHref(string(body), url)
		zap.S().Debugf("Processed %s found %d links", url, len(hrefs))

		pq.AddURLs(hrefs, scraperItem)
		scraper.WriteFile(body, url)
	}
}
