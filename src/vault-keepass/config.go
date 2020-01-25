package main

import "net/http"

// Configuration - configuration
type Configuration struct {
	VaultURL         string
	VaultToken       string
	VaultPath        string
	VaultInsecureSSL bool
	VaultTimeout     int
}

// VaultKVResult - Result from Vault GET request
type VaultKVResult struct {
	RequestID string            `json:"request_id"`
	Data      map[string]string `json:"data"`
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
