package MetricsGetter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestK8sMetricsGetter_record(t *testing.T) {
	testGetter := NewK8sMetricsGetter(make(chan bool), make(chan string), make(chan bool))

	res := testGetter.recordMetrics()
	fmt.Println(res)
}

func TestK8sMetricsGetter_Start(t *testing.T) {
	testGetter := NewK8sMetricsGetter(make(chan bool), make(chan string), make(chan bool))

	go func() {
		assert.Equal(t, true, <-testGetter.controlCh)
		assert.Equal(t, "test_file_name", <-testGetter.fileNameCh)
		testGetter.retCh <- true
	}()

	res := testGetter.Start("test_file_name")
	assert.Equal(t, true, res)
}

func TestK8sMetricsGetter_Stop(t *testing.T) {
	testGetter := NewK8sMetricsGetter(make(chan bool), make(chan string), make(chan bool))

	go func() {
		assert.Equal(t, false, <-testGetter.controlCh)
		testGetter.retCh <- true
	}()

	res := testGetter.Stop()
	assert.Equal(t, true, res)
}

func TestK8sMetricsGetter_Run(t *testing.T) {
	testGetter := NewK8sMetricsGetter(make(chan bool), make(chan string), make(chan bool))

	go testGetter.Run()

	res := testGetter.Start("test_file_name")
	assert.Equal(t, true, res)
	res = testGetter.Stop()
	assert.Equal(t, true, res)
}
