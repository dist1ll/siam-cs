package main

import (
	siam "github.com/m2q/algo-siam"
	"github.com/m2q/algo-siam/client"
	"github.com/m2q/siam-cs"
	"log"
)

func main() {

	// Create AlgorandBuffer
	b, err := siam.CreateAlgorandBuffer(client.CreateAlgorandClientMock("", ""), "")
	if err != nil {
		log.Fatal(err)
	}

	// Start Oracle
	err = csgo.NewOracle(b, &csgo.HLTV{}).Serve()
	if err != nil {
		log.Fatal(err)
	}
}
