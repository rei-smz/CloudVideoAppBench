package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/Bench"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		log.Println("Please specify your test config file")
		return
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
		return
	}
	bench := Bench.NewClientBench(os.Args[1])
	if bench == nil {
		log.Fatalln("Failed initialising bench.")
		return
	}

	bench.Bench()
}
