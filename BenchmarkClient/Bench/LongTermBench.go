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

type LongTermBenchUser struct {
	User
	requestIDs []string
}

func NewLongTermBenchUser(benchConfig *BenchConfig) IUser {
	return &LongTermBenchUser{
		User: User{
			httpController: HTTPController.NewHTTPController(),
			benchConfig:    benchConfig,
		},
	}
}

func (lb *LongTermBenchUser) onSendingRequests() {
	// Send http request
	objectPath := lb.benchConfig.PathPrefix +
		strconv.Itoa(rand.Intn(lb.benchConfig.UserIdRange))
	reqBody := maps.Clone(lb.benchConfig.ReqArgs)
	reqBody["path"] = objectPath
	reqBody["type"] = "request"
	response := lb.httpController.Post(lb.benchConfig.URL, reqBody)
	//log.Println(response)
	if response != nil && response["status_code"] == 200 {
		lb.requestIDs = append(lb.requestIDs, response["key"].(string))
	}
	// Wait
	time.Sleep(time.Duration(lb.benchConfig.UserWaiting) * time.Second)
}

func (lb *LongTermBenchUser) onStopSendingRequest() {
	//log.Println(lb.requestIDs)
	time.Sleep(time.Minute * 5)
	reqBody := make(map[string]any)
	reqBody["type"] = "query"
	reqBody["args"] = make(map[string]string)
	for _, id := range lb.requestIDs {
		log.Println("Checking request " + id)
		reqBody["args"].(map[string]string)["key"] = id
		for {
			response := lb.httpController.Post(lb.benchConfig.URL, reqBody)
			if response != nil && response["status_code"] == 200 {
				if response["status"] == "running" {
					time.Sleep(5 * time.Minute)
				} else {
					result := make(map[string]any)
					result["result"] = response["status"]
					result["time_stamp"] = int64(response["start_time"].(float64))
					result["rtt"] = time.Duration((response["end_time"].(float64) - response["start_time"].(float64)) * 1e6).Milliseconds()
					result["status_code"] = response["status_code"]
					lb.results = append(lb.results, result)
					break
				}
			}
		}
	}
}

// Request sends user requests and handles responses every benchConfig.UserWaiting seconds.
func (lb *LongTermBenchUser) Request(stopCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

	for {
		select {
		case <-stopCh:
			lb.onStopSendingRequest()
			return
		default:
			lb.onSendingRequests()
		}
	}
}
