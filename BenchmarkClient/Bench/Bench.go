package Bench

import (
	"fmt"
	"log"
	"main/RPCClient"
	"os"
	"sync"
	"time"
)

type ClientBench struct {
	benchConfig *BenchConfig
	//resultsMu   sync.Mutex
	Users     []IUser
	rpcClient RPCClient.RPCClient
}

// NewClientBench creates a LongTermBenchUser object.
func NewClientBench(confPath string) *ClientBench {
	config := LoadConfig(confPath)
	if config == nil {
		log.Fatalln("Failed loading benchConfig.")
		return nil
	}

	rpcClient := RPCClient.NewRPCClient()
	if rpcClient == nil {
		log.Fatalln("Failed initialising rpcClient")
		return nil
	}

	return &ClientBench{
		benchConfig: config,
		rpcClient:   rpcClient,
	}
}

func (cb *ClientBench) Bench() {
	// Send start recording request to server
	err := cb.rpcClient.Connect()
	if err != nil {
		log.Fatalln("Failed to connect to server module")
		return
	}
	err = cb.rpcClient.RequestStart(cb.benchConfig.TestName + time.Now().String())
	if err != nil {
		log.Println("Failed to start recording on server module. You may need to start it manually.")
	}

	log.Println("Long term benchmark is running")

	// Create goroutine for each user
	stopCh := make(chan struct{})
	var wg sync.WaitGroup
	for range cb.benchConfig.NumUser {
		var newUser IUser
		switch cb.benchConfig.AppType {
		case "long_term":
			newUser = NewLongTermBenchUser(cb.benchConfig)
		case "default":
			newUser = NewDefaultBenchUser(cb.benchConfig)
		default:
			newUser = NewDefaultBenchUser(cb.benchConfig)
		}
		cb.Users = append(cb.Users, newUser)
		wg.Add(1)
		go newUser.Request(stopCh, &wg)
	}

	// Wait for duration end and turn off goroutines
	time.Sleep(time.Duration(cb.benchConfig.Duration) * time.Second)
	close(stopCh)
	wg.Wait()

	// Send stop recording request to server
	err = cb.rpcClient.RequestStop()
	if err != nil {
		log.Println("Failed to stop recording on server module. You may need to stop it manually.")
	}
	cb.rpcClient.CloseConn()

	cb.saveResults()
}

// saveResults saves results to a local file.
func (cb *ClientBench) saveResults() {
	totalRequests := 0
	errorRequests := 0
	var totalRTT float64 = 0

	file, err := os.Create(cb.benchConfig.TestName + "-" + time.Now().Format("2006-01-02-15-04-05") + ".txt")
	if err != nil {
		log.Fatalln("Can not create result file")
		return
	}
	defer file.Close()

	// Read results and process
	file.WriteString("time_stamp\tstatus_code\tresponse_time\n")

	for _, user := range cb.Users {
		userResults := user.GetResults()
		totalRequests += len(userResults)
		for _, res := range userResults {
			ws := fmt.Sprintf("%v\t%v\t%v\n", res["time_stamp"], res["status_code"], res["rtt"])
			file.WriteString(ws)
			if res["status_code"] != 200 || res["result"] == "error" {
				errorRequests++
			} else {
				totalRTT += float64(res["rtt"].(int64))
			}
		}
	}

	errorRate := float64(errorRequests) / float64(totalRequests)
	avgRTT := totalRTT / (float64(totalRequests))
	// Save summary
	file.WriteString("===============\n")
	file.WriteString(fmt.Sprintf("Total requests: %v, error rate: %v, average RTT: %v",
		totalRequests, errorRate, avgRTT))
}
