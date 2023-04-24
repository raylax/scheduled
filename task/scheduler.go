package task

import "time"

type s struct {
	skip bool
}

type Scheduler struct {
	tasks      []*Task
	q          *Q
	now        NowFunc
	shutdownCh chan struct{}
	ticker     *time.Ticker
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		q:     newQ(),
		tasks: make([]*Task, 0),
		now: func() int64 {
			return time.Now().Unix()
		},
		ticker:     time.NewTicker(1 * time.Second),
		shutdownCh: make(chan struct{}),
	}
}

func (s *Scheduler) AddTask(task *Task) {
	s.tasks = append(s.tasks, task)
	ts := task.ticker.next(s.now())
	s.q.Push(ts, task)
}

func (s *Scheduler) Start() error {
	for {

		select {
		case t := <-s.ticker.C:
			ok, ts, tasks := s.q.Peek()
			if !ok || ts > t.Unix() {
				continue
			}
			s.q.Pop()
			for _, task := range tasks {
				next := task.ticker.next(s.now())
				s.q.Push(next, task)
				go task.f()
			}
		case <-s.shutdownCh:
			return nil
		}

	}
}

func (s *Scheduler) Shutdown() {
	s.shutdownCh <- struct{}{}
}
