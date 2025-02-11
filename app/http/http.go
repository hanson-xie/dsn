package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func LoadSqlConfig(url, tomlPath string) (*http.Response, error) {
	var client = &http.Client{
		Timeout: 30 * time.Second,
	}
	reqBody := map[string]string{
		"toml_path": tomlPath,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	return client.Do(req)
}
