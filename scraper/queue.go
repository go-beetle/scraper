package scraper

import (
	"container/heap"
	"go.uber.org/zap"
	"time"
)

// An Url is something we manage in a ScrapeTime queue.
type Url struct {
	Url                   string // The Url of the item; arbitrary.
	ScrapeTime            int    // The ScrapeTime of the item in the queue.
	ScrapeIntervalSeconds int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Url

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].ScrapeTime < pq[j].ScrapeTime
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Url)
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

// Sleep until the current time is less than the time on the prioritized item
func (pq *PriorityQueue) GetNextUrl() string {
	last := (*pq)[0]
	timeNow := int(time.Now().Unix())

	if timeNow < last.ScrapeTime {
		t := last.ScrapeTime - timeNow
		zap.S().Debugf("Sleeping for %d seconds", t)
		time.Sleep(time.Duration(t) * time.Second)
	}

	last.ScrapeTime = last.ScrapeIntervalSeconds + timeNow
	heap.Fix(pq, 0)
	return last.Url
}