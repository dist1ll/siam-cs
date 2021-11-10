package csgo

import (
	"testing"

	siam "github.com/m2q/algo-siam"
	"github.com/m2q/algo-siam/client"
)

// setupOracleMockedAPI creates an Oracle instance, initialized with a
// mocked AlgorandBuffer and mocked API.
func setupOracleMockedAPI() *Oracle {
	c := client.CreateAlgorandClientMock("", "")
	buffer, _ := siam.CreateAlgorandBuffer(c, "")
	cfg := &OracleConfig{
		PrimaryAPI:      nil,
		RefreshInterval: 0,
	}
	o := NewOracle(buffer, cfg)
	return o
}

func TestOracle(t *testing.T) {
	oracle := setupOracleMockedAPI()

	oracle.Serve()
}
