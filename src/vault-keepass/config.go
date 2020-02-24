package main

import "net/http"

// Configuration - configuration
type Configuration struct {
	VaultURL         string
	VaultToken       string
	VaultPath        string
	VaultInsecureSSL bool
	VaultTimeout     int
	Debug            bool
}

// VaultKVResult - Result from Vault GET request
type VaultKVResult struct {
	RequestID string            `json:"request_id"`
	Data      map[string]string `json:"data"`
}

// VaultKVListResult - result from LIST request
type VaultKVListResult struct {
	RequestID string       `json:"request_id"`
	Data      VaultKeyList `json:"data"`
}

// VaultKeyList - list ov keys
type VaultKeyList struct {
	Keys []string `json:"keys"`
}

// VaultError - Vault error format
type VaultError struct {
	Messages []string
}

// HTTPResult - Result of HTTP operation
type HTTPResult struct {
	Status     string
	StatusCode int
	Header     http.Header
	Content    []byte
}
