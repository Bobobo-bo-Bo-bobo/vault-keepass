package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func cmdList(cfg *Configuration, args []string) error {
	var res VaultKVListResult
	var resp HTTPResult
	var err error

	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Error: list command don't accept additional parameters\n\n")
		showUsage()
		os.Exit(1)
	}

	resp, err = httpRequest(cfg, cfg.VaultURL+cfg.VaultPath, "LIST", nil, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			// Note: LIST does only work on directories - see https://www.vaultproject.io/api-docs/secret/kv/kv-v1/
			//       so we try GET to see if it is a file instead

			resp, err = httpRequest(cfg, cfg.VaultURL+cfg.VaultPath, "GET", nil, nil)
			if err != nil {
				return err
			}

			if resp.StatusCode != http.StatusOK {
				log.WithFields(log.Fields{
					"http_status":         resp.StatusCode,
					"http_status_message": resp.Status,
				}).Error("Invalid HTTP status received")
				return fmt.Errorf("Listing of Vault keys failed")
			}

		} else {
			log.WithFields(log.Fields{
				"http_status":         resp.StatusCode,
				"http_status_message": resp.Status,
			}).Error("Invalid HTTP status received")
			return fmt.Errorf("Listing of Vault keys failed")
		}
	}

	err = json.Unmarshal(resp.Content, &res)
	if err != nil {
		return err
	}
	for _, k := range res.Data.Keys {
		fmt.Println(k)
	}

	return nil
}
