package MetricsGetter

import (
	"log"
	"os"
)

type MetricsGetter interface {
	Run()
	Start(fileName string) bool
	Stop() bool
}

// NewMetricsGetter is a factory function creating an object of the concrete classes implementing interface MetricsGetter
func NewMetricsGetter() MetricsGetter {
	metricsType := os.Getenv("METRICS_TYPE")
	if metricsType == "" {
		log.Println("METRICS_TYPE missing")
		return nil
	}

	switch metricsType {
	case "k8s":
		return NewK8sMetricsGetter(make(chan bool), make(chan string), make(chan bool))
	default:
		return nil
	}
}
