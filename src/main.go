package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	"os"
)

var redisConnectionPool *redis.Pool

func main() {

	simulate := flag.Bool("sim", false, "simulate instead of starting a web snake")
	setWeightsFlag := flag.Bool("set", false, "set weights for snake")
	flag.Parse()
	redisConnectionPool = NewPool()

	if *setWeightsFlag {
		weights := map[string]int{}
		weights["hug-walls"] = 0
		weights["straight"] = 0
		weights["random"] = 0
		weights["control"] = 100
		weights["nearest-food"] = 100
		weights["agressive"] = 100
		weights["attempt-kill"] = 100
		weights["avoid-death"] = 100

		StoreWeights(weights)
		fmt.Printf("Wrote: %v", weights)
		return
	}
	if !*simulate {
		http.HandleFunc("/start", start)
		http.HandleFunc("/move", move)
		port := os.Getenv("PORT")
		if port == "" {
			port = "9000"
		}
		PrimeWeightsCache()
		log.Printf("Running server on port %s...\n", port)
		http.ListenAndServe(":"+port, nil)
	} else {
		log.Println("Simulate a game to train swunetics!")
		numWorkers := 10
		numGames := 200
		numFood := 6
		bestQuality := TrainAgainstSnek(numGames, 0, numFood, numWorkers, 0)
		fmt.Printf("\nstarting quality: %v\n", bestQuality)
		for {
			quality := TrainAgainstSnek(numGames, 5, numFood, numWorkers, bestQuality)
			if quality > bestQuality {
				bestQuality = quality
			}
		}
	}
}
