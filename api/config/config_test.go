// config/config_test.go

package config_test

import (
	"encoding/hex"
	"github.com/avkos/file-registry/api/config"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Success_WithSetEnv(t *testing.T) {
	// Reset the global Config before the test
	config.Config = config.GlobalConfig{}
	// Set environment variables programmatically
	t.Setenv("CONTRACT_ADDRESS", "0x0000000000000000000000000000000000000002")
	t.Setenv("ETH_RPC_URL", "http://localhost:8546")
	t.Setenv("IPFS_URL", "http://localhost:5002")
	t.Setenv("PORT", "8001")
	t.Setenv("CHAIN_ID", "1338")
	t.Setenv("PRIVATE_KEY", "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	err := config.LoadConfig()
	assert.NoError(t, err, "Expected no error with environment variables set via t.Setenv")

	//// Verify Config fields
	expectedAddress := common.HexToAddress("0x0000000000000000000000000000000000000002")
	assert.Equal(t, expectedAddress, config.Config.ContractAddress, "ContractAddress mismatch")

	assert.Equal(t, "http://localhost:8546", config.Config.EthRpcUrl, "EthRpcUrl mismatch")
	assert.Equal(t, "http://localhost:5002", config.Config.IpfsUrl, "IpfsUrl mismatch")
	assert.Equal(t, "8001", config.Config.Port, "Port mismatch")

	expectedChainID := big.NewInt(1338)
	assert.Equal(t, expectedChainID, config.Config.ChainID, "ChainID mismatch")

	expectedPrivateKey, _ := hex.DecodeString("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	assert.Equal(t, expectedPrivateKey, config.Config.PrivateKey, "PrivateKey mismatch")
}

func TestLoadConfig_MissingRequiredVariables(t *testing.T) {
	// Reset the global Config before the test
	config.Config = config.GlobalConfig{}
	t.Setenv("CONTRACT_ADDRESS", "0x0000000000000000000000000000000000000002")
	t.Setenv("ETH_RPC_URL", "")
	t.Setenv("IPFS_URL", "")
	t.Setenv("PORT", "")
	t.Setenv("CHAIN_ID", "")
	t.Setenv("PRIVATE_KEY", "")

	err := config.LoadConfig()
	assert.Error(t, err, "Expected validation error due to missing required variables")
	assert.Contains(t, err.Error(), "config validation error", "Error message should indicate validation issues")
}
