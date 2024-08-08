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
	metricsGetter MetricsGetter.MetricsGetter
}

func (c *requestController) Start(ctx context.Context, in *pb.StartRequest) (*pb.ServerReply, error) {
	log.Println("Received start request")
	return &pb.ServerReply{Success: c.metricsGetter.Start(in.FileName)}, nil
}

func (c *requestController) Stop(ctx context.Context, in *pb.StopRequest) (*pb.ServerReply, error) {
	log.Println("Received stop request")
	return &pb.ServerReply{Success: c.metricsGetter.Stop()}, nil
}

type RPCServer struct {
	grpcServer *grpc.Server
}

func NewRPCServer(metricsGetter MetricsGetter.MetricsGetter) *RPCServer {
	s := grpc.NewServer()
	reqController := requestController{
		metricsGetter: metricsGetter,
	}
	pb.RegisterRecodingControlServer(s, &reqController)

	return &RPCServer{
		grpcServer: s,
	}
}

func (s *RPCServer) ListenAndServe() {
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
