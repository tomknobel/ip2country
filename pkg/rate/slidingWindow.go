package rate

import (
	"sync"
	"time"
)

type Window interface {
	Start() time.Time
	Count() int64
	AddCount(n int64)
	Reset(s time.Time, c int64)
}

// LocalWindow is a static window the previous window and the current one
type LocalWindow struct {
	start int64
	count int64
}

func NewLocalWindow() *LocalWindow {
	return &LocalWindow{}
}

func (w *LocalWindow) Start() time.Time {
	return time.Unix(0, w.start)
}

func (w *LocalWindow) Count() int64 {
	return w.count
}

func (w *LocalWindow) AddCount(n int64) {
	w.count += n
}

func (w *LocalWindow) Reset(s time.Time, c int64) {
	w.start = s.UnixNano()
	w.count = c
}

type Limiter struct {
	size  time.Duration
	limit int64

	mu sync.Mutex

	curr Window
	prev Window
}

func NewLimiter(size time.Duration, limit int64) *Limiter {
	currWin := NewLocalWindow()

	prevWin := NewLocalWindow()

	lim := &Limiter{
		size:  size,
		limit: limit,
		curr:  currWin,
		prev:  prevWin,
	}

	return lim
}

func (lim *Limiter) Allow() bool {
	return lim.allowN(time.Now(), 1)
}

func (lim *Limiter) advance(now time.Time) {
	newCurrStart := now.Truncate(lim.size)

	diffSize := newCurrStart.Sub(lim.curr.Start()) / lim.size
	if diffSize >= 1 {
		newPrevCount := int64(0)
		if diffSize == 1 {
			newPrevCount = lim.curr.Count()
		}
		lim.prev.Reset(newCurrStart.Add(-lim.size), newPrevCount)

		lim.curr.Reset(newCurrStart, 0)
	}
}

// allowN reports whether n events may happen at time now.
func (lim *Limiter) allowN(now time.Time, n int64) bool {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	lim.advance(now)

	elapsed := now.Sub(lim.curr.Start())
	weight := float64(lim.size-elapsed) / float64(lim.size)
	count := int64(weight*float64(lim.prev.Count())) + lim.curr.Count()
	if count+n > lim.limit {
		return false
	}

	lim.curr.AddCount(n)
	return true
}
