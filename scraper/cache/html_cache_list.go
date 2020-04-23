package cache

type cacheHeap []*HTMLCacheItem

func (pq cacheHeap) Len() int { return len(pq) }

func (pq cacheHeap) Less(i, j int) bool {
	return pq[i].scrapeFrequency < pq[j].scrapeFrequency
}

func (pq cacheHeap) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *cacheHeap) Push(x interface{}) {
	*pq = append(*pq, x.(*HTMLCacheItem))
}

func (pq *cacheHeap) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	*pq = old[0 : n-1]
	return item
}