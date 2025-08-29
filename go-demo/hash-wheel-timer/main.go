package hashwheeltimer

import (
	"container/list"
	"sync"
	"time"
)

// 任务
type TimerTask struct {
	expiration      int64
	remainingRounds int
	fn              func()
	element         *list.Element
}

// 时间轮的每个桶结构
type Bucket struct {
	tasks *list.List
	mu    sync.Mutex
}

func NewBucket() *Bucket {
	return &Bucket{
		tasks: list.New(),
	}
}

func (b *Bucket) Add(task *TimerTask) {
	b.mu.Lock()
	defer b.mu.Unlock()
	task.element = b.tasks.PushBack(task)
}

func (b *Bucket) CollectExpired(now int64) []*TimerTask {
	b.mu.Lock()
	defer b.mu.Unlock()
	var due []*TimerTask
	for e := b.tasks.Front(); e != nil; {
		task := e.Value.(*TimerTask)
		next := e.Next()
		if task.remainingRounds > 0 {
			task.remainingRounds--
			e = next
			continue
		}
		if task.expiration <= now {
			b.tasks.Remove(e)
			task.element = nil
			due = append(due, task)
		}
		e = next
	}
	return due
}

// 任务执行线程池
type Executor struct {
	workers int
	queue   chan func()
	wg      sync.WaitGroup
	quit    chan struct{}
}

func NewExecutor(workers, queueCap int) *Executor {
	return &Executor{
		workers: workers,
		queue:   make(chan func(), queueCap),
		quit:    make(chan struct{}),
	}
}

func (e *Executor) Start() {
	for i := 0; i < e.workers; i++ {
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			for {
				select {
				case job := <-e.queue:
					safeRun(job)
				case <-e.quit:
					return
				}
			}
		}()
	}
}

func (e *Executor) Stop() {
	close(e.quit)
	e.wg.Wait()
}

// Submit 将任务投递到线程池；若队列满，为避免阻塞 tick，降级为 go 执行
func (e *Executor) Submit(job func()) {
	select {
	case e.queue <- job:
	default:
		// 降级兜底策略：避免阻塞 tick（也可以改为丢弃或阻塞）
		go safeRun(job)
	}
}

func safeRun(job func()) {
	defer func() {
		if r := recover(); r != nil {
			// 避免业务 panic 杀死 worker
		}
	}()
	job()
}

// 时间轮
type HashedWheelTimer struct {
	tick       time.Duration
	wheelSize  int
	wheel      []*Bucket
	currentPos int
	excutor    *Executor
	stopCh     chan struct{}
	startOnce  sync.Once
}

func NewHashedWheelTimer(tick time.Duration, wheelSize int, excutor *Executor) *HashedWheelTimer {
	wheel := make([]*Bucket, wheelSize)
	for i := range wheel {
		wheel[i] = NewBucket()
	}
	return &HashedWheelTimer{
		tick:      tick,
		wheelSize: wheelSize,
		wheel:     wheel,
		excutor:   excutor,
		stopCh:    make(chan struct{}),
	}
}

func (h *HashedWheelTimer) Start() {
	h.startOnce.Do(func() {
		h.excutor.Start()
	})
	go h.worker()
}

func (h *HashedWheelTimer) Stop() {
	close(h.stopCh)
	if h.excutor != nil {
		h.excutor.Stop()
	}
}

func (h *HashedWheelTimer) worker() {
	ticker := time.NewTicker(h.tick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.currentPos = (h.currentPos + 1) % h.wheelSize
			now := time.Now().UnixNano()
			bucket := h.wheel[h.currentPos]

			// 把时间轮中到期的定时任务推入执行器内的goroutine执行
			due := bucket.CollectExpired(now)
			for _, t := range due {
				h.excutor.Submit(t.fn)
			}

		case <-h.stopCh:
			return
		}
	}
}

func (h *HashedWheelTimer) NewTimeout(delay time.Duration, fn func()) {
	if delay < 0 {
		delay = 0
	}
	expiration := time.Now().Add(delay).UnixNano()
	ticks := int(delay / h.tick)
	remainingRounds := ticks / h.wheelSize
	pos := (h.currentPos + ticks) % h.wheelSize

	task := &TimerTask{
		fn:              fn,
		expiration:      expiration,
		remainingRounds: remainingRounds,
	}

	h.wheel[pos].Add(task)
}
