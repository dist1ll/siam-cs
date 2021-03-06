package main

import (
	"log"
	"time"

	siam "github.com/m2q/algo-siam"
	"github.com/m2q/siam-cs"
)

func main() {

	// Create AlgorandBuffer
	b, err := siam.NewAlgorandBufferFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Configure Oracle
	cfg := &csgo.OracleConfig{
		PrimaryAPI:      &csgo.HLTV{},
		RefreshInterval: time.Minute * 3,
	}

	// Create Oracle
	oracle := csgo.NewOracle(b, cfg)

	// Start Oracle
	oracle.Serve()

	// Wait until oracle finishes
	oracle.Wait()
}
