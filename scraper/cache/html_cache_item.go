package cache

type HTMLCacheItem struct {
	url             string
	body            []byte
	scrapeFrequency int `nanoseconds`
	size            int
}
