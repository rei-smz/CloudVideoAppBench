package RPCClient

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "main/RecordingControl"
	"os"
)

type RPCClient interface {
	Connect() error
	CloseConn()
	RequestStart(fileName string) error
	RequestStop() error
}

type rpcClient struct {
	rpcURI string
	conn   *grpc.ClientConn
	client pb.RecodingControlClient
}

func NewRPCClient() RPCClient {
	rpcURI := os.Getenv("RPC_URI")
	if rpcURI == "" {
		log.Println("No RPC_URI found")
		return nil
	}

	return &rpcClient{
		rpcURI: rpcURI,
	}
}

func (c *rpcClient) Connect() error {
	conn, err := grpc.NewClient(c.rpcURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		conn.Close()
		return err
	}
	client := pb.NewRecodingControlClient(conn)

	c.conn = conn
	c.client = client
	return nil
}

func (c *rpcClient) CloseConn() {
	c.conn.Close()
}

func (c *rpcClient) RequestStart(fileName string) error {
	retVal, err := c.client.Start(context.Background(), &pb.StartRequest{FileName: fileName})
	if err != nil {
		return err
	}

	if retVal.Success == false {
		return errors.New("recording did not start successfully")
	}

	return nil
}

func (c *rpcClient) RequestStop() error {
	retVal, err := c.client.Stop(context.Background(), &pb.StopRequest{})
	if err != nil {
		return err
	}

	if retVal.Success == false {
		return errors.New("recording did not start successfully")
	}

	return nil
}
