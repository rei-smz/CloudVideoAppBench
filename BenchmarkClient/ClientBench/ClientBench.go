package ClientBench

import (
	"fmt"
	"log"
	"main/HTTPController"
	"main/RPCClient"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type ClientBench struct {
	properties *BenchProperties
	results    []map[string]any
	resultsMu  sync.Mutex
	rpcClient  RPCClient.RPCClient
}

// NewClientBench creates a ClientBench object.
func NewClientBench() *ClientBench {
	properties := LoadProperties()
	if properties == nil {
		log.Fatalln("Failed loading properties.")
		return nil
	}

	rpcClient := RPCClient.NewRPCClient()
	if rpcClient == nil {
		log.Fatalln("Failed initialising rpcClient")
		return nil
	}

	return &ClientBench{
		properties: properties,
		rpcClient:  rpcClient,
	}
}

func (cb *ClientBench) Bench() {
	// Send start recording request to server
	err := cb.rpcClient.Connect()
	if err != nil {
		log.Fatalln("Failed to connect to server module")
		return
	}
	err = cb.rpcClient.RequestStart(cb.properties.TestName + time.Now().String())
	if err != nil {
		log.Println("Failed to start recording on server module. You may start it manually.")
	}

	// Create goroutine for each user
	stopCh := make(chan struct{})
	var wg sync.WaitGroup
	for range cb.properties.NumUser {
		wg.Add(1)
		go cb.request(HTTPController.NewHTTPController(), stopCh, &wg)
	}

	// Wait for duration end and turn off goroutines
	time.Sleep(time.Duration(cb.properties.Duration) * time.Millisecond)
	close(stopCh)
	wg.Wait()

	// Send stop recording request to server
	err = cb.rpcClient.RequestStop()
	if err != nil {
		log.Println("Failed to stop recording on server module. You may stop it manually.")
	}
	cb.rpcClient.CloseConn()

	cb.saveResults()
}

// request sends user requests and handle responses every properties.UserWaiting seconds.
func (cb *ClientBench) request(httpController HTTPController.HTTPController, stopCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	//time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

	for {
		select {
		case <-stopCh:
			return
		default:
			// Send http request
			objectPath := cb.properties.PathPrefix +
				strconv.Itoa(rand.Intn(cb.properties.UserIdRange)) +
				"/" + cb.properties.ObjectName
			reqBody := make(map[string]string)
			reqBody["path"] = objectPath
			startT := time.Now()
			response := httpController.Post(cb.properties.URL, reqBody)
			timeCost := time.Since(startT)
			// Handle response
			if response == nil {
				log.Println("Client module error")
			} else {
				result := response
				result["rtt"] = timeCost.Milliseconds()
				result["time_stamp"] = startT.UnixMilli()
				cb.resultsMu.Lock()
				cb.results = append(cb.results, result)
				cb.resultsMu.Unlock()
			}
			// Wait
			time.Sleep(time.Duration(cb.properties.UserWaiting) * time.Millisecond)
		}
	}
}

// saveResults saves results to a local file.
func (cb *ClientBench) saveResults() {
	totalRequests := len(cb.results)
	errorRequests := 0
	totalRTT := 0

	file, err := os.OpenFile(cb.properties.TestName+"-"+time.Now().String(), os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		log.Fatalln("Can not create result file")
		return
	}
	defer file.Close()

	// Read results and process
	file.WriteString("time stamp\tstatus code\tresponse time\n")
	for _, res := range cb.results {
		ws := fmt.Sprintf("%v\t%v\t%v\n", res["time_stamp"], res["status_code"].(string), res["rtt"].(int))
		file.WriteString(ws)
		if res["status_code"].(string) != "200" {
			errorRequests++
		}
		totalRTT += res["rtt"].(int)
	}

	errorRate := float64(errorRequests) / float64(totalRequests)
	avgRTT := float64(totalRTT) / (float64(totalRequests - errorRequests))
	// Save summary
	file.WriteString("===============\n")
	file.WriteString(fmt.Sprintf("Total requests: %v, error rate: %v, average RTT: %v",
		totalRequests, errorRate, avgRTT))
}
