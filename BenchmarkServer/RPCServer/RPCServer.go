package RPCServer

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"main/MetricsGetter"
	pb "main/RecordingControl"
	"net"
)

type requestController struct {
	pb.UnimplementedRecodingControlServer
	controlCh   chan bool
	fileNameCh  chan string
	getterRetCh chan bool
}

func (c *requestController) Start(ctx context.Context, in *pb.StartRequest) (*pb.ServerReply, error) {
	log.Println("Received start request")
	c.controlCh <- true
	c.fileNameCh <- in.FileName
	return &pb.ServerReply{Success: <-c.getterRetCh}, nil
}

func (c *requestController) Stop(ctx context.Context, in *pb.StopRequest) (*pb.ServerReply, error) {
	log.Println("Received stop request")
	c.controlCh <- false
	return &pb.ServerReply{Success: <-c.getterRetCh}, nil
}

type RPCServer struct {
	grpcServer    *grpc.Server
	metricsGetter *MetricsGetter.MetricsGetter
}

func NewRPCServer() *RPCServer {
	s := grpc.NewServer()
	reqController := requestController{
		controlCh:   make(chan bool),
		fileNameCh:  make(chan string),
		getterRetCh: make(chan bool),
	}
	pb.RegisterRecodingControlServer(s, &reqController)

	return &RPCServer{
		grpcServer: s,
		metricsGetter: MetricsGetter.NewMetricsGetter(reqController.controlCh,
			reqController.fileNameCh, reqController.getterRetCh),
	}
}

func (s *RPCServer) ListenAndServe() {
	go s.metricsGetter.Run()

	log.Println("Server is starting...")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
		return
	}
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
		return
	}
}
