package contracts_test

import (
	"github.com/avkos/file-registry/api/config"
	"testing"

	"github.com/avkos/file-registry/api/contracts"
	"github.com/stretchr/testify/assert"
)

func TestLoadTransactor_InvalidKey(t *testing.T) {
	auth, err := contracts.LoadTransactor()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse private key")
	assert.Nil(t, auth)
}

func TestLoadTransactor_Success(t *testing.T) {
	t.Setenv("CONTRACT_ADDRESS", "0x0000000000000000000000000000000000000002")
	t.Setenv("ETH_RPC_URL", "http://localhost:8546")
	t.Setenv("IPFS_URL", "http://localhost:5002")
	t.Setenv("PORT", "8001")
	t.Setenv("CHAIN_ID", "1338")
	t.Setenv("PRIVATE_KEY", "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	err := config.LoadConfig()
	assert.NoError(t, err)

	auth, err := contracts.LoadTransactor()
	assert.NoError(t, err)
	assert.NotNil(t, auth)
}
