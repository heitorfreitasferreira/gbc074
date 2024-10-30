package database

import "time"

const invalidationTime time.Duration = 5 * time.Second

const (
	Cluster0Replica0Path string = "/tmp/grupo9/user/0"
	Cluster0Replica1Path string = "/tmp/grupo9/user/1"
	Cluster0Replica2Path string = "/tmp/grupo9/user/2"
	Cluster1Replica0Path string = "/tmp/grupo9/book/0"
	Cluster1Replica1Path string = "/tmp/grupo9/book/1"
	Cluster1Replica2Path string = "/tmp/grupo9/book/2"
)

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
