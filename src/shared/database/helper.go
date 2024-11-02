package database

import "time"

const invalidationTime time.Duration = 5 * time.Second

type Timed[T any] struct {
	created time.Time
	Item    T
}

func NewTimed[T any](item T) Timed[T] {
	return Timed[T]{created: time.Now(), Item: item}
}

func (t Timed[T]) IsValid() bool {
	return time.Since(t.created) < invalidationTime
}
