package ClientBench

import (
	"log"
	"os"
	"strconv"
)

type BenchProperties struct {
	NumUser     int
	Duration    int
	UserIdRange int
	UserWaiting int
	URL         string
	PathPrefix  string
	ObjectName  string
	TestName    string
}

func LoadProperties() *BenchProperties {
	numUser, err := loadValue("NUM_USER")
	if err != nil {
		return nil
	}
	duration, err := loadValue("DURATION")
	if err != nil {
		return nil
	}
	userIdRange, err := loadValue("USER_ID_RANGE")
	if err != nil {
		return nil
	}
	url := os.Getenv("URL")
	if url == "" {
		log.Println("URL missing")
		return nil
	}
	userWaiting, err := loadValue("USER_WAITING")
	if err != nil {
		return nil
	}
	objectName := os.Getenv("OBJECT")
	if objectName == "" {
		log.Println("OBJECT missing")
		return nil
	}
	testName := os.Getenv("TEST_NAME")
	if testName == "" {
		log.Println("TEST_NAME missing")
		return nil
	}
	pathPre := os.Getenv("PATH_PREFIX")
	if pathPre == "" {
		log.Println("PATH_PREFIX missing")
		return nil
	}

	return &BenchProperties{
		NumUser:     numUser,
		URL:         url,
		Duration:    duration,
		UserIdRange: userIdRange,
		UserWaiting: userWaiting,
		ObjectName:  objectName,
		TestName:    testName,
		PathPrefix:  pathPre,
	}
}

func loadValue(key string) (int, error) {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		log.Printf("%s missing or invalid\n", key)
	}
	return value, err
}
