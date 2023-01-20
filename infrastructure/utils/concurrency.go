package utils

import "sync"

type Result[T any] struct {
	Data  T
	Error error
}

func WaitToCloseChan[T any](wg *sync.WaitGroup, ch chan T) {
	wg.Wait()
	close(ch)
}
