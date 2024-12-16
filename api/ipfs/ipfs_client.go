package ipfs

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/ipfs/boxo/files"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/ipfs/kubo/client/rpc"
)

type IPFSClient struct {
	apiURL string
	api    *rpc.HttpApi
}

type RPCUnixfsAPI interface {
	Unixfs() RPCUnixfs
}

type RPCUnixfs interface {
	Add(context.Context, files.File) (path.Resolved, error)
}

// NewIPFSClient creates a new IPFS client using NewApiWithClient, connecting to the given IPFS API URL.
func NewIPFSClient(apiURL string) (*IPFSClient, error) {
	api, err := rpc.NewURLApiWithClient(apiURL, http.DefaultClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create IPFS API with client: %w", err)
	}

	return &IPFSClient{
		apiURL: apiURL,
		api:    api,
	}, nil
}

// Add uploads the given file content to IPFS using the Unixfs API and returns the CID.
func (c *IPFSClient) Add(ctx *gin.Context, fileContent []byte) (string, error) {
	f := files.NewBytesFile(fileContent)

	p, err := c.api.Unixfs().Add(ctx, f)
	if err != nil {
		return "", fmt.Errorf("failed to add file to IPFS: %w", err)
	}
	return p.RootCid().String(), nil
}
