package main

import (
	siam "github.com/m2q/algo-siam"
	"github.com/m2q/algo-siam/client"
	"github.com/m2q/siam-cs"
	"log"
	"time"
)

func main() {

	// Create AlgorandBuffer
	b, err := siam.CreateAlgorandBuffer(client.CreateAlgorandClientMock("", ""), "")
	if err != nil {
		log.Fatal(err)
	}

	// Configure Oracle
	cfg := &csgo.OracleConfig{
		PrimaryAPI:      &csgo.HLTV{},
		RefreshInterval: time.Minute,
	}

	// Start Oracle
	err = csgo.NewOracle(b, cfg).Serve()
	if err != nil {
		log.Fatal(err)
	}
}
