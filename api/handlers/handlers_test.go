package handlers_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/avkos/file-registry/api/handlers"
)

// Mock implementations for Contract and IPFSClient
type mockContract struct {
	saveFunc func(filePath, cid string) (string, error)
	getFunc  func(filePath string) (string, error)
}

func (m *mockContract) Save(filePath, cid string) (string, error) {
	return m.saveFunc(filePath, cid)
}

func (m *mockContract) Get(filePath string) (string, error) {
	return m.getFunc(filePath)
}

type mockIPFSClient struct {
	addFunc func(ctx *gin.Context, file []byte) (string, error)
}

func (m *mockIPFSClient) Add(ctx *gin.Context, file []byte) (string, error) {
	return m.addFunc(ctx, file)
}

// TestUploadFile_Success tests that uploading a file returns a 200 status and the expected CID.
func TestUploadFile_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mock contract and IPFS
	mockC := &mockContract{
		saveFunc: func(filePath, cid string) (string, error) {
			// Return a fake txHash
			return "0x1234567890abcdef", nil
		},
	}
	mockIPFS := &mockIPFSClient{
		addFunc: func(ctx *gin.Context, file []byte) (string, error) {
			return "QmFakeCID", nil
		},
	}

	router := handlers.SetupRouter(mockC, mockIPFS)

	// Prepare a valid base64 file
	fileContent := []byte("Hello World!")
	fileB64 := base64.StdEncoding.EncodeToString(fileContent)
	body := []byte(`{"filePath":"/test/file.txt","file":"` + fileB64 + `"}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/files", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "QmFakeCID", resp["cid"])
	assert.Equal(t, "0x1234567890abcdef", resp["txHash"])
}

// TestUploadFile_InvalidBase64 tests that invalid base64 data returns a 400 status.
func TestUploadFile_InvalidBase64(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockC := &mockContract{}
	mockIPFS := &mockIPFSClient{}

	router := handlers.SetupRouter(mockC, mockIPFS)

	body := []byte(`{"filePath":"/test/file.txt","file":"!!!notbase64!!!"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/files", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "Invalid base64 data")
}

// TestUploadFile_ContractError tests that if the contract save fails, we return 500.
func TestUploadFile_ContractError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockC := &mockContract{
		saveFunc: func(filePath, cid string) (string, error) {
			return "", errors.New("contract save failed")
		},
	}
	mockIPFS := &mockIPFSClient{
		addFunc: func(ctx *gin.Context, file []byte) (string, error) {
			return "QmFakeCID", nil
		},
	}

	router := handlers.SetupRouter(mockC, mockIPFS)

	fileContent := []byte("Hello World!")
	fileB64 := base64.StdEncoding.EncodeToString(fileContent)
	body := []byte(`{"filePath":"/test/file.txt","file":"` + fileB64 + `"}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/files", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "contract save failed")
}

func TestUploadFile_MissingFilePath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockC := &mockContract{}
	mockIPFS := &mockIPFSClient{}

	router := handlers.SetupRouter(mockC, mockIPFS)

	fileContent := []byte("Hello World!")
	fileB64 := base64.StdEncoding.EncodeToString(fileContent)
	body := []byte(`{"file":"` + fileB64 + `"}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/files", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "Missing filePath query parameter")
}

// TestGetFile_Success tests that retrieving a file CID returns 200 and correct CID.
func TestGetFile_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockC := &mockContract{
		getFunc: func(filePath string) (string, error) {
			return "QmFakeCID", nil
		},
	}
	mockIPFS := &mockIPFSClient{}

	router := handlers.SetupRouter(mockC, mockIPFS)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/files?filePath=/test/file.txt", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "QmFakeCID", resp["cid"])
}

// TestGetFile_MissingFilePath tests that if filePath is not provided, returns 400.
func TestGetFile_MissingFilePath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockC := &mockContract{}
	mockIPFS := &mockIPFSClient{}

	router := handlers.SetupRouter(mockC, mockIPFS)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/files", nil)
	fmt.Println(req)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "Missing filePath query parameter")
}

// TestGetFile_ContractError tests that if contract get fails, returns 500.
func TestGetFile_ContractError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockC := &mockContract{
		getFunc: func(filePath string) (string, error) {
			return "", errors.New("contract get failed")
		},
	}
	mockIPFS := &mockIPFSClient{}

	router := handlers.SetupRouter(mockC, mockIPFS)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/files?filePath=/test/file.txt", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "contract get failed")
}
