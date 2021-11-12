package csgo

import (
	siam "github.com/m2q/algo-siam"
	"github.com/m2q/algo-siam/client"
	"github.com/m2q/siam-cs/generator"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// setupOracleMockedAPI creates an Oracle instance, initialized with a
// AlgorandBuffer (client.AlgorandMock) and StubAPI.
func setupOracleMockedAPI() (*Oracle, *siam.AlgorandBuffer, *StubAPI) {
	c := client.CreateAlgorandClientMock("", "")
	buffer, err := siam.CreateAlgorandBuffer(c, client.GeneratePrivateKey64())
	if err != nil {
		panic(err)
	}
	api := &StubAPI{}
	cfg := &OracleConfig{
		PrimaryAPI:      api,
		RefreshInterval: 0,
		SiamCfg:         &siam.ManageConfig{},
	}
	o := NewOracle(buffer, cfg)
	return o, buffer, api
}

// Tests if the Oracle writes the data provided by the API stub to the AlgorandBuffer
func TestOracle_SimpleSmallData(t *testing.T) {
	oracle, buffer, stub := setupOracleMockedAPI()
	oracle.Serve()

	// Set match data to stub
	past, future := generator.GetData()
	stub.SetMatches(past, future)

	// Check if desired state is written by Oracle
	desired := ConstructDesiredState(past, future, client.GlobalBytes)
	contains := buffer.ContainsWithin(CreateWinnerMap(desired), time.Second*5)
	assert.True(t, contains)

	// Stop oracle
	oracle.Stop()
}
