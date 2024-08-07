package MetricsGetter

import (
	"fmt"
	"testing"
)

func TestMetricsGetter_record(t *testing.T) {
	testGetter := NewMetricsGetter(make(chan bool), make(chan string), make(chan bool))

	res := testGetter.recordMetrics()
	fmt.Println(res)
}
