package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/ClientBench"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return
	}
	bench := ClientBench.NewClientBench()
	if bench == nil {
		log.Fatalln("Failed initialising bench.")
		return
	}

	bench.Bench()
}
