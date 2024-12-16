package handlers

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Contract interface {
	Save(filePath string, cid string) (string, error)
	Get(filePath string) (string, error)
}

type IPFSClient interface {
	Add(ctx *gin.Context, file []byte) (string, error)
}

type Handlers struct {
	Contract   Contract
	IPFSClient IPFSClient
}

type FileUploadRequest struct {
	FilePath string `json:"filePath"`
	FileB64  string `json:"file"`
}

func (h *Handlers) UploadFile(c *gin.Context) {
	var req FileUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON: " + err.Error()})
		return
	}
	if req.FilePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing filePath query parameter"})
		return
	}
	fileBytes, err := base64.StdEncoding.DecodeString(req.FileB64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base64 data: " + err.Error()})
		return
	}

	cid, err := h.IPFSClient.Add(c, fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "IPFS Add error: " + err.Error()})
		return
	}

	txHash, err := h.Contract.Save(req.FilePath, cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Contract save error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cid": cid, "txHash": txHash})
}

func (h *Handlers) GetFile(c *gin.Context) {
	filePath := c.Query("filePath")

	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing filePath query parameter"})
		return
	}

	cid, err := h.Contract.Get(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Contract get error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cid": cid})
}

func SetupRouter(contract Contract, ipfsClient IPFSClient) *gin.Engine {
	h := &Handlers{Contract: contract, IPFSClient: ipfsClient}
	router := gin.Default()
	router.POST("/v1/files", h.UploadFile)
	router.GET("/v1/files", h.GetFile)
	return router
}
