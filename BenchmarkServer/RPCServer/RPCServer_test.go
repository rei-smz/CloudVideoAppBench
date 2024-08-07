package RPCServer

import (
	"context"
	"github.com/stretchr/testify/assert"
	pb "main/RecordingControl"
	"testing"
)

func TestRequestController_Start(t *testing.T) {
	testControlCh := make(chan bool)
	testFileNameCh := make(chan string)
	testGetterRetCh := make(chan bool)
	testRC := &requestController{
		controlCh:   testControlCh,
		fileNameCh:  testFileNameCh,
		getterRetCh: testGetterRetCh,
	}

	go func() {
		assert.Equal(t, true, <-testControlCh)
		assert.Equal(t, "test_file_name", <-testFileNameCh)
		testGetterRetCh <- true
	}()
	reply, err := testRC.Start(context.Background(), &pb.StartRequest{FileName: "test_file_name"})

	assert.NoError(t, err)
	assert.Equal(t, true, reply.Success)
}

func TestRequestController_Stop(t *testing.T) {
	testControlCh := make(chan bool)
	testFileNameCh := make(chan string)
	testGetterRetCh := make(chan bool)
	testRC := &requestController{
		controlCh:   testControlCh,
		fileNameCh:  testFileNameCh,
		getterRetCh: testGetterRetCh,
	}

	go func() {
		assert.Equal(t, false, <-testControlCh)
		testGetterRetCh <- true
	}()
	reply, err := testRC.Stop(context.Background(), &pb.StopRequest{})

	assert.NoError(t, err)
	assert.Equal(t, true, reply.Success)
}
