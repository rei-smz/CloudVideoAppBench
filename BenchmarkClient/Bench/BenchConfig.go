package Bench

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type BenchConfig struct {
	AppType     string         `json:"app_type"`
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
}
