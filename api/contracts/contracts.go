package contracts

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"os"
	"path/filepath"

	"github.com/avkos/file-registry/api/config"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type FileRegistry struct {
	FileRegistryCaller
	FileRegistryTransactor
}

type FileRegistryCaller struct {
	contract *bind.BoundContract
}

type FileRegistryTransactor struct {
	contract *bind.BoundContract
}

// LoadTransactor loads a transactor using the PRIVATE_KEY from config.
func LoadTransactor() (*bind.TransactOpts, error) {
	privateKey, err := crypto.ToECDSA(config.Config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, config.Config.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor with chain ID: %w", err)
	}
	return auth, nil
}

// NewFileRegistry creates a new instance of FileRegistry, bound to a specific deployed contract.
func NewFileRegistry(address common.Address, backend *ethclient.Client) (*FileRegistry, error) {
	abiJson, err := loadABI(filepath.Join("contracts", "file_registry.abi"))
	if err != nil {
		return nil, err
	}
	contract := bind.NewBoundContract(address, abiJson, backend, backend, backend)
	return &FileRegistry{
		FileRegistryCaller:     FileRegistryCaller{contract: contract},
		FileRegistryTransactor: FileRegistryTransactor{contract: contract},
	}, nil
}

func loadABI(filePath string) (abi.ABI, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return abi.ABI{}, err
	}
	return abi.JSON(bytes.NewReader(file))
}

func (f *FileRegistryTransactor) Save(opts *bind.TransactOpts, filePath string, cid string) (*types.Transaction, error) {
	return f.contract.Transact(opts, "save", filePath, cid)
}

func (f *FileRegistryCaller) Get(opts *bind.CallOpts, filePath string) (string, error) {
	var out []interface{}
	err := f.contract.Call(opts, &out, "get", filePath)
	if err != nil {
		return "", err
	}
	return out[0].(string), nil

}

// ContractAPI provides a simpler interface that handlers can use directly.
// It wraps the FileRegistry contract and a TransactOpts for sending transactions.
type ContractAPI struct {
	instance *FileRegistry
	auth     *bind.TransactOpts
	client   *ethclient.Client
}

// NewContractAPI connects to the Ethereum client, loads the contract and transactor.
func NewContractAPI() (*ContractAPI, error) {
	// Connect to Ethereum
	client, err := ethclient.Dial(config.Config.EthRpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum: %w", err)
	}

	// Load contract instance
	registry, err := NewFileRegistry(config.Config.ContractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate contract: %w", err)
	}

	// Load transactor
	auth, err := LoadTransactor()
	if err != nil {
		return nil, fmt.Errorf("failed to load transactor: %w", err)
	}

	return &ContractAPI{
		instance: registry,
		auth:     auth,
		client:   client,
	}, nil
}

// Save stores the CID for the given filePath on-chain.
func (api *ContractAPI) Save(filePath, cid string) (string, error) {
	tx, err := api.instance.Save(api.auth, filePath, cid)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// Get retrieves the CID for the given filePath from the contract.
func (api *ContractAPI) Get(filePath string) (string, error) {
	return api.instance.Get(nil, filePath)
}
