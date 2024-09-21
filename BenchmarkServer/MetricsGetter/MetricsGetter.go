package MetricsGetter

import (
	"log"
	"os"
	"time"
)

type MetricsGetter interface {
	Run()
	Start(fileName string) bool
	Stop() bool
	recordMetrics() string
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

func (g *metricsGetter) Run() {
	log.Println("MetricsGetter started")
	running := false
	var file *os.File

	for {
		select {
		case ctrl := <-g.controlCh:
			if ctrl {
				if file != nil {
					file.Close()
				}
				var err error
				file, err = os.Create(<-g.fileNameCh)
				if err != nil {
					file.Close()
					g.retCh <- false
				} else {
					g.retCh <- true
					running = true
					log.Println("Metrics getter is running")
				}
			} else {
				err := file.Close()
				running = false
				log.Println("Metrics getter is stopped")
				if err != nil {
					g.retCh <- false
				} else {
					file = nil
					g.retCh <- true
				}
			}
		default:
			if running {
				file.WriteString(g.recordMetrics())
				time.Sleep(15 * time.Second)
			}
		}
	}
}

func (g *metricsGetter) recordMetrics() string {
	return ""
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
