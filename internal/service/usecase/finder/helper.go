package finder

type IntervalHeap []*Interval

func (h IntervalHeap) Len() int {
	return len(h)
}
func (h IntervalHeap) Less(i, j int) bool {
	return h[i].From.Before(h[j].From)
}
func (h IntervalHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntervalHeap) Push(x interface{}) {
	*h = append(*h, x.(*Interval))
}

func (h *IntervalHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *IntervalHeap) Peek() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	return x
}
