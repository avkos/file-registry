package e2e

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var apiBaseURL = "http://localhost:8000/v1"

type uploadResp struct {
	CID    string `json:"cid"`
	TxHash string `json:"txHash"`
}

type getResp struct {
	CID string `json:"cid"`
}

func TestE2E_UploadAndRetrieveFile(t *testing.T) {
	waitForService(t, apiBaseURL+"/files")

	fileContent := []byte("Hello Boxo E2E Test!")
	fileB64 := base64.StdEncoding.EncodeToString(fileContent)
	filePath := "/e2e/boxo_test_file.txt"

	// Upload file
	uploadBody := map[string]interface{}{
		"filePath": filePath,
		"file":     fileB64,
	}
	ur := uploadFile(t, uploadBody)
	require.NotEmpty(t, ur.CID)
	require.NotEmpty(t, ur.TxHash)

	cid := getFileCID(t, filePath)
	assert.Equal(t, ur.CID, cid)
}

func uploadFile(t *testing.T, body map[string]interface{}) uploadResp {
	b, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, apiBaseURL+"/files", bytes.NewReader(b))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var ur uploadResp
	err = json.NewDecoder(resp.Body).Decode(&ur)
	require.NoError(t, err)
	return ur
}

func getFileCID(t *testing.T, filePath string) string {
	req, err := http.NewRequest(http.MethodGet, apiBaseURL+"/files?filePath="+filePath, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var gr getResp
	err = json.NewDecoder(resp.Body).Decode(&gr)
	require.NoError(t, err)

	return gr.CID
}

func waitForService(t *testing.T, url string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode < 500 {
			resp.Body.Close()
			break
		}
		if ctx.Err() != nil {
			t.Fatalf("Service not ready after 30s: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
}
