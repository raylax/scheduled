package task

import (
	"container/heap"
	"sync"
)

type Q struct {
	q *q
	m map[int64]*qi
	l sync.Mutex
}

func newQ() *Q {
	q := &q{}
	heap.Init(q)
	return &Q{q: q, m: map[int64]*qi{}}
}

func (q *Q) IsEmpty() bool {
	q.l.Lock()
	defer q.l.Unlock()
	if q.q.Len() == 0 {
		return true
	}
	return false
}

func (q *Q) Push(ts int64, task *Task) {
	q.l.Lock()
	defer q.l.Unlock()
	if i, ok := q.m[ts]; ok {
		i.tasks = append(i.tasks, task)
		return
	}
	heap.Push(q.q, &qi{
		ts:    ts,
		tasks: []*Task{task},
	})
}

func (q *Q) Peek() (bool, int64, []*Task) {
	q.l.Lock()
	defer q.l.Unlock()
	if q.q.Len() == 0 {
		return false, 0, nil
	}
	item := q.q.top().(*qi)
	return true, item.ts, item.tasks
}

func (q *Q) Pop() (bool, int64, []*Task) {
	q.l.Lock()
	defer q.l.Unlock()
	if q.q.Len() == 0 {
		return false, 0, nil
	}
	item := heap.Pop(q.q).(*qi)
	return true, item.ts, item.tasks
}

type qi struct {
	ts    int64
	tasks []*Task
}

type q []*qi

func (q q) top() any {
	return q[0]
}

func (q q) Len() int {
	return len(q)
}

func (q q) Less(i, j int) bool {
	return q[i].ts > q[j].ts
}

func (q q) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *q) Push(x any) {
	item := x.(*qi)
	*q = append(*q, item)
}

func (q *q) Pop() any {
	old := *q
	n := len(old)
	item := old[n-1]
	*q = old[0 : n-1]
	return item
}
