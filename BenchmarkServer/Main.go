package main

import "main/RPCServer"

func main() {
	rpcServer := RPCServer.NewRPCServer()

	rpcServer.ListenAndServe()
}
