package Bench

import (
	"main/HTTPController"
	"sync"
)

type User struct {
	results        []map[string]any
	httpController HTTPController.HTTPController
	benchConfig    *BenchConfig
}

type IUser interface {
	Request(stopCh <-chan struct{}, wg *sync.WaitGroup)
	onSendingRequests()
	onStopSendingRequest()
	GetResults() []map[string]any
}

func (u *User) GetResults() []map[string]any {
	return u.results
}
