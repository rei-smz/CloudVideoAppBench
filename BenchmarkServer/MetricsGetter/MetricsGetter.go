package MetricsGetter

import (
	"log"
)

type MetricsGetter interface {
	Run()
	Start(fileName string) bool
	Stop() bool
}

type metricsGetter struct {
	controlCh  chan bool
	fileNameCh chan string
	retCh      chan bool
}

func (g *metricsGetter) Start(fileName string) bool {
	g.controlCh <- true
	g.fileNameCh <- fileName
	return <-g.retCh
}

func (g *metricsGetter) Stop() bool {
	g.controlCh <- false
	return <-g.retCh
}

// NewMetricsGetter is a factory function creating an object of the concrete classes implementing interface metricsGetter
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
