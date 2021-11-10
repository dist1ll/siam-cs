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
	wg, _ := csgo.NewOracle(b, cfg).Serve()

	// Wait until oracle finishes
	wg.Wait()
}
