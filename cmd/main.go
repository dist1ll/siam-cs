package main

import (
	"log"
	"time"

	siam "github.com/m2q/algo-siam"
	"github.com/m2q/algo-siam/client"
	"github.com/m2q/siam-cs"
)

func main() {
	// Create AlgorandBuffer
	b, err := siam.CreateAlgorandBuffer(client.CreateAlgorandClientMock("", ""),
		client.GeneratePrivateKey64())
	if err != nil {
		log.Fatal(err)
	}

	// Configure Oracle
	cfg := &csgo.OracleConfig{
		PrimaryAPI:      &csgo.HLTV{},
		RefreshInterval: time.Minute,
	}

	// Create Oracle
	oracle := csgo.NewOracle(b, cfg)

	// Start Oracle
	oracle.Serve()

	// Wait until oracle finishes
	oracle.Wait()
}
