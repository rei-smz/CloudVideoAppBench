package MetricsGetter

import (
	"log"
)

type MetricsGetter interface {
	Run()
	Start(fileName string) bool
	Stop() bool
}

// NewMetricsGetter is a factory function creating an object of the concrete classes implementing interface MetricsGetter
func NewMetricsGetter(metricsType string) MetricsGetter {
	//metricsType := os.Getenv("METRICS_TYPE")

	switch metricsType {
	case "k8s":
		{
			return NewK8sMetricsGetter(make(chan bool), make(chan string), make(chan bool))
		}
	default:
		{
			log.Println("METRICS_TYPE missing or invalid")
			return nil
		}
	}
}
