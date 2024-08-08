package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/MetricsGetter"
	"main/RPCServer"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return
	}

	metricsGetter := MetricsGetter.NewMetricsGetter()
	if metricsGetter == nil {
		log.Fatalln("Cannot create metrics getter")
		return
	}
	rpcServer := RPCServer.NewRPCServer(metricsGetter)

	go metricsGetter.Run()
	rpcServer.ListenAndServe()
}
