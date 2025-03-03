package Bench

import (
	"log"
	"main/HTTPController"
	"maps"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type DefaultBenchUser struct {
	User
}

func NewDefaultBenchUser(benchConfig *BenchConfig) IUser {
	return &DefaultBenchUser{
		User: User{
			httpController: HTTPController.NewHTTPController(),
			benchConfig:    benchConfig,
		},
	}
}

func (db *DefaultBenchUser) onSendingRequests() {
	// Send http request
	objectPath := db.benchConfig.PathPrefix +
		strconv.Itoa(rand.Intn(db.benchConfig.UserIdRange))
	reqBody := maps.Clone(db.benchConfig.ReqArgs)
	reqBody["path"] = objectPath

	startT := time.Now()
	response := db.httpController.Post(db.benchConfig.URL, reqBody)
	timeCost := time.Since(startT)
	// Handle response
	if response == nil {
		log.Println("Client module error")
	} else {
		result := response
		result["rtt"] = timeCost.Milliseconds()
		result["time_stamp"] = startT.UnixMilli()
		db.results = append(db.results, result)
	}
	// Wait
	time.Sleep(time.Duration(db.benchConfig.UserWaiting) * time.Second)
}

// Request sends user requests and handles responses every benchConfig.UserWaiting seconds.
func (db *DefaultBenchUser) Request(stopCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

	for {
		select {
		case <-stopCh:
			db.onStopSendingRequest()
			return
		default:
			db.onSendingRequests()
		}
	}
}

func (db *DefaultBenchUser) onStopSendingRequest() {}
