package RPCServer

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"main/MetricsGetter"
	pb "main/RecordingControl"
	"testing"
)

func TestRequestController_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMetricsGetter := MetricsGetter.NewMockMetricsGetter(ctrl)
	testRC := &requestController{
		metricsGetter: mockMetricsGetter,
	}

	mockMetricsGetter.EXPECT().Start("test_file_name").Return(true)
	reply, err := testRC.Start(context.Background(), &pb.StartRequest{FileName: "test_file_name"})

	assert.NoError(t, err)
	assert.Equal(t, true, reply.Success)
}

func TestRequestController_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMetricsGetter := MetricsGetter.NewMockMetricsGetter(ctrl)
	testRC := &requestController{
		metricsGetter: mockMetricsGetter,
	}

	mockMetricsGetter.EXPECT().Stop().Return(true)
	reply, err := testRC.Stop(context.Background(), &pb.StopRequest{})

	assert.NoError(t, err)
	assert.Equal(t, true, reply.Success)
}
