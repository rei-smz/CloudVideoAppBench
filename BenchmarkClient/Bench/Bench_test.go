package Bench

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"main/HTTPController"
	"main/RPCClient"
	"sync"
	"testing"
)

func NewTestClient(client RPCClient.RPCClient) *LongTermBenchUser {
	return &LongTermBenchUser{
		benchConfig: &BenchConfig{
			NumUser:     1,
			Duration:    45,
			UserIdRange: 1,
			UserWaiting: 5,
			URL:         "testURL",
			PathPrefix:  "testPath",
			ObjectName:  "testObj",
			TestName:    "test",
			ReqArgs:     make(map[string]string),
		},
		rpcClient: client,
	}
}

func TestClientBench_Bench(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRpc := RPCClient.NewMockRPCClient(ctrl)
	testClient := NewTestClient(mockRpc)
	mockRpc.EXPECT().Connect().Return(nil).Times(1)
	mockRpc.EXPECT().RequestStart(gomock.Any()).Return(nil).Times(1)
	mockRpc.EXPECT().RequestStop().Return(nil).Times(1)
	testClient.Bench()
}

func TestClientBench_Request(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRpc := RPCClient.NewMockRPCClient(ctrl)
	testClient := NewTestClient(mockRpc)
	mockHttp := HTTPController.NewMockHTTPController(ctrl)
	mockHttp.EXPECT().Post(gomock.Any(), gomock.Any()).Return(map[string]any{"status_code": 200}).AnyTimes()
	stopCh := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go testClient.request(mockHttp, stopCh, &wg)
	close(stopCh)
	wg.Wait()
	for _, res := range testClient.results {
		assert.Equal(t, 200, res["status_code"])
	}
}
