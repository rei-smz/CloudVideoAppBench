package RPCClient

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Copy .env file to this folder before testing

func TestRpcClient_Connect(t *testing.T) {
	godotenv.Load()
	testRpc := NewRPCClient()
	err := testRpc.Connect()
	defer testRpc.CloseConn()
	assert.NoError(t, err)
}

func TestRpcClient_RequestStart(t *testing.T) {
	godotenv.Load()
	testRpc := NewRPCClient()
	testRpc.Connect()
	defer testRpc.CloseConn()

	err := testRpc.RequestStart("test_file_name_unit_test")
	assert.NoError(t, err)
}

func TestRpcClient_RequestStop(t *testing.T) {
	godotenv.Load()
	testRpc := NewRPCClient()
	testRpc.Connect()
	defer testRpc.CloseConn()

	err := testRpc.RequestStop()
	assert.NoError(t, err)
}
