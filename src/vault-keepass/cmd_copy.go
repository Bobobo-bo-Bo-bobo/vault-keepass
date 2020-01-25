package main

import (
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func cmdCopy(cfg *Configuration, args []string) error {
	var res VaultKVResult
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Error: One and only one key must be given\n\n")
		showUsage()
		os.Exit(1)
	}

	key := args[1]

	resp, err := httpRequest(cfg, cfg.VaultURL+cfg.VaultPath, "GET", nil, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"http_status":         resp.StatusCode,
			"http_status_message": resp.Status,
		}).Error("Invalid HTTP status received")
		return fmt.Errorf("Getting list of Vault keys failed")
	}

	err = json.Unmarshal(resp.Content, &res)
	if err != nil {
		return err
	}

	v, found := res.Data[key]
	if !found {
		return fmt.Errorf("Key not found")
	}

	err = clipboard.WriteAll(v)
	if err != nil {
		return err
	}

	return nil
}
