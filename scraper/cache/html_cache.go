package cache

import "container/heap"

type HTMLCache struct {
	itemMap map[string]*HTMLCacheItem
	heap cacheHeap
	maxSize int
	currentSize int
}

func InitHTMLCache(size int) HTMLCache {
	return HTMLCache{
		itemMap:     make(map[string]*HTMLCacheItem),
		heap:        make([]*HTMLCacheItem, 0),
		maxSize:     size,
		currentSize: 0,
	}
}

// Returns nil if not found
func (c *HTMLCache) Get(url string) []byte {
	if item, ok := c.itemMap[url]; ok {
		return item.body
	}
	return nil
}

func (c *HTMLCache) Add(url string, body []byte, freq int) error {
	len_body := len(body)

	if len_body > c.maxSize {
		return nil // Change this to error
	}

	// Add to heap
	c.currentSize += len_body
	item := &HTMLCacheItem{
		url:             url,
		body:            body,
		scrapeFrequency: freq,
		size:            len_body,
	}
	heap.Push(&c.heap, item)
	c.itemMap[url] = item

	for c.currentSize > c.maxSize {
		item := heap.Pop(&c.heap).(*HTMLCacheItem)

		c.currentSize -= item.size
		delete(c.itemMap, item.url)
	}
	return nil
}