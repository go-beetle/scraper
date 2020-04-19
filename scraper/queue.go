package scraper

import (
	"container/heap"
	"go.uber.org/zap"
	"time"
)

// An ScraperItem is something we manage in a ScrapeTime queue.
type ScraperItem struct {
	Url                       *URL // The ScraperItem of the item; arbitrary.
	ScrapeTime                int64    // The ScrapeTime of the item in the queue.
	ScrapeIntervalNanoseconds int64
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*ScraperItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].ScrapeTime < pq[j].ScrapeTime
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*ScraperItem)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func InitPQSeed(urls []*ScraperItem) PriorityQueue {
	pq := make(PriorityQueue, len(urls))
	now := time.Now().UnixNano()
	for i, u := range urls {
		pq[i] = &ScraperItem{
			Url:                       u.Url,
			ScrapeTime:                now,
			ScrapeIntervalNanoseconds: u.ScrapeIntervalNanoseconds,
		}
	}
	heap.Init(&pq)
	return pq
}

// Sleep until the current time is less than the time on the prioritized item
func (pq *PriorityQueue) PeekScraperItemAndUpdate() *ScraperItem {
	last := (*pq)[0]
	timeNow := time.Now().UnixNano()

	if timeNow < last.ScrapeTime {
		t := last.ScrapeTime - timeNow
		zap.S().Debugf("Sleeping for %dms", t / 1e6)
		time.Sleep(time.Duration(t) * time.Nanosecond)
	}

	timeNow = time.Now().UnixNano()
	last.ScrapeTime = last.ScrapeIntervalNanoseconds + timeNow
	heap.Fix(pq, 0)
	return last
}

// Adds a newly discovered url to priority queue
// Copies over specified scrape interval from the URL that
// contained this link
func (pq *PriorityQueue) AddURLs(urls []*URL, from *ScraperItem) {
	timeNow := time.Now().UnixNano()
	for _, u := range urls {
		item := &ScraperItem{
			Url:                       u,
			ScrapeTime:                timeNow,
			ScrapeIntervalNanoseconds: from.ScrapeIntervalNanoseconds,
		}
		heap.Push(pq, item)
	}
}