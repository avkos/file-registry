package config

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"math/big"
	"strings"
)

type Validation struct {
	ContractAddress string `envconfig:"CONTRACT_ADDRESS" validate:"required,len=42,startswith=0x"`
	EthRpcUrl       string `envconfig:"ETH_RPC_URL" validate:"required,url"`
	IpfsUrl         string `envconfig:"IPFS_URL" validate:"required,url"`
	Port            string `envconfig:"PORT" validate:"required,numeric"`
	ChainID         string `envconfig:"CHAIN_ID" validate:"numeric"`
	PrivateKeyHex   string `envconfig:"PRIVATE_KEY" validate:"required,hexadecimal,len=66,startswith=0x"`
}

type GlobalConfig struct {
	ContractAddress common.Address
	EthRpcUrl       string
	IpfsUrl         string
	Port            string
	ChainID         *big.Int
	PrivateKey      []byte
}

var Config GlobalConfig

// LoadConfig loads environment variables into Config, then validates them.
func LoadConfig() error {
	// Load from .env if present (not required)
	err := godotenv.Load(".env")
	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		fmt.Println("failed to load .env file. Continuing without it.")
	}
	var cfg Validation
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("failed to load config from env: %w", err)
	}

	// Validate config using validator
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return fmt.Errorf("config validation error: %w", err)
	}

	// Parse and assign to global Config struct
	Config.ContractAddress = common.HexToAddress(cfg.ContractAddress)
	if Config.ContractAddress == (common.Address{}) {
		return fmt.Errorf("invalid CONTRACT_ADDRESS: %s", cfg.ContractAddress)
	}

	Config.EthRpcUrl = cfg.EthRpcUrl
	Config.IpfsUrl = cfg.IpfsUrl
	Config.Port = cfg.Port

	Config.ChainID = big.NewInt(1) // default to 1 if CHAIN_ID not provided
	if cfg.ChainID != "" {
		if _, ok := Config.ChainID.SetString(cfg.ChainID, 10); !ok {
			return fmt.Errorf("invalid CHAIN_ID: %s", cfg.ChainID)
		}
	}

	// Remove 0x prefix from PRIVATE_KEY if present
	privateKeyHex := strings.TrimPrefix(cfg.PrivateKeyHex, "0x")
	if privateKeyHex == "" {
		return fmt.Errorf("PRIVATE_KEY not set or invalid")
	}
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to decode private key hex: %w", err)
	}
	Config.PrivateKey = privateKeyBytes

	return nil
}
