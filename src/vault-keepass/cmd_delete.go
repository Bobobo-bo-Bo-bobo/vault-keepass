package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

func cmdDelete(cfg *Configuration, args []string) error {
	var keys []string

	if len(args) == 1 {
		fmt.Fprintf(os.Stderr, "Error: Key(s) to delete are required for this command\n\n")
		showUsage()
		os.Exit(1)
	}

	// remove duplicates
	temp := make(map[string]bool)
	for _, k := range args[1:] {
		temp[k] = true
	}
	for k := range temp {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		resp, err := httpRequest(cfg, cfg.VaultURL+filepath.Join(cfg.VaultPath, k), "DELETE", nil, nil)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
			log.WithFields(log.Fields{
				"http_status":         resp.StatusCode,
				"http_status_message": resp.Status,
				"key":                 k,
			}).Error("Invalid HTTP status received")
			return fmt.Errorf("Removal of Vault key failed")
		}
	}

	return nil
}
