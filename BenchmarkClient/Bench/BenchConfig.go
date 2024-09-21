package Bench

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type BenchConfig struct {
	NumUser     int            `json:"n_user"`
	Duration    int            `json:"duration"`
	UserIdRange int            `json:"id_range"`
	UserWaiting int            `json:"wait"`
	URL         string         `json:"url"`
	PathPrefix  string         `json:"path_prefix"`
	TestName    string         `json:"test_name"`
	ReqArgs     map[string]any `json:"req_args"`
}

func LoadConfig(confPath string) *BenchConfig {
	jsonFile, err := os.Open(confPath)
	if err != nil {
		log.Fatalln("Config file missing")
		return nil
	}
	defer jsonFile.Close()

	byteVal, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln("Unable to open config file")
		return nil
	}

	var config BenchConfig
	err = json.Unmarshal(byteVal, &config)
	if err != nil {
		log.Fatalln("Unable to unmarshal, check your json file")
		return nil
	}

	return &config

	//numUser, err := loadValue("NUM_USER")
	//if err != nil {
	//	return nil
	//}
	//duration, err := loadValue("DURATION")
	//if err != nil {
	//	return nil
	//}
	//userIdRange, err := loadValue("USER_ID_RANGE")
	//if err != nil {
	//	return nil
	//}
	//url := os.Getenv("URL")
	//if url == "" {
	//	log.Println("URL missing")
	//	return nil
	//}
	//userWaiting, err := loadValue("USER_WAITING")
	//if err != nil {
	//	return nil
	//}
	//objectName := os.Getenv("OBJECT")
	//if objectName == "" {
	//	log.Println("OBJECT missing")
	//	return nil
	//}
	//testName := os.Getenv("TEST_NAME")
	//if testName == "" {
	//	log.Println("TEST_NAME missing")
	//	return nil
	//}
	//reqArgs := os.Getenv("REQ_ARGS")
	//if reqArgs == "" {
	//	log.Println("REQ_ARGS missing")
	//	return nil
	//}
	//
	//return &BenchConfig{
	//	NumUser:     numUser,
	//	URL:         url,
	//	Duration:    duration,
	//	UserIdRange: userIdRange,
	//	UserWaiting: userWaiting,
	//	ObjectName:  objectName,
	//	TestName:    testName,
	//	ReqArgs:     reqArgs,
	//}
}

//func loadValue(key string) (int, error) {
//	value, err := strconv.Atoi(os.Getenv(key))
//	if err != nil {
//		log.Printf("%s missing or invalid\n", key)
//	}
//	return value, err
//}
