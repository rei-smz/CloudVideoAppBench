package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/MetricsGetter"
	"main/RPCServer"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return
	}

	metricsGetter := MetricsGetter.NewMetricsGetter(os.Getenv("METRICS_TYPE"))
	if metricsGetter == nil {
		log.Fatalln("Cannot create metrics getter")
		return
	}
	rpcServer := RPCServer.NewRPCServer(metricsGetter)

	go metricsGetter.Run()
	rpcServer.ListenAndServe()
}
