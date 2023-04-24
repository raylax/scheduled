package task

import (
	cronexpr "github.com/raylax/scheduled/cron"
	"time"
)

type Type int

const (
	TypeNone Type = iota
	TypeCron
)

type Task struct {
	id     string
	t      Type
	ticker ticker
	f      func()
}

func NewTask(id string, cron string, f func()) (*Task, error) {
	expression, err := cronexpr.Parse(cron)
	if err != nil {
		return nil, err
	}

	return &Task{
		id:     id,
		t:      TypeCron,
		f:      f,
		ticker: &cronTicker{expression: expression},
	}, nil
}

type ticker interface {
	next(now int64) int64
}

type cronTicker struct {
	expression *cronexpr.Expression
}

func (c *cronTicker) next(now int64) int64 {
	return c.expression.Next(time.Unix(now, 0)).Unix()
}
