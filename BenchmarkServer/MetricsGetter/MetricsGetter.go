package MetricsGetter

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type MetricsGetter struct {
	controlCh  chan bool
	fileNameCh chan string
	retCh      chan bool
}

func NewMetricsGetter(ctrlCh chan bool, fnCh chan string, retCh chan bool) *MetricsGetter {
	return &MetricsGetter{
		controlCh:  ctrlCh,
		fileNameCh: fnCh,
		retCh:      retCh,
	}
}

func (g *MetricsGetter) Run() {
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

func (g *MetricsGetter) recordMetrics() string {
	timestamp := time.Now().UnixMilli()

	cmd := exec.Command("kubectl", "top", "-n", "openfaas-fn", "pod")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintln(timestamp, "Error executing kubectl:", err)
	}

	return fmt.Sprintln(timestamp, string(output))
}
