package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func cmdSet(cfg *Configuration, args []string) error {
	var existing = make(map[string]string)
	var res VaultKVResult
	var payload []byte
	var err error
	var key string
	var value string

	parse := flag.NewFlagSet("cmd-set", flag.ExitOnError)
	var replace = parse.Bool("replace", false, "Replace data instead of merging")
	parse.Parse(args[1:])
	parse.Usage = showUsage

	_args := parse.Args()

	if len(_args) != 2 && len(_args) != 0 {
		fmt.Fprintf(os.Stderr, "Error: Not enough arguments for 'set' command\n\n")
		showUsage()
		os.Exit(1)
	}

	if len(_args) == 0 {
		// Read standard, don't display input
		fmt.Print("Key: ")
		key = readStandardInput()
		fmt.Println()
		fmt.Print("Value: ")
		value = readStandardInput()
	} else {
		key = _args[0]
		value = _args[1]
	}

	if !*replace {
		// fetch existing values
		resp, err := httpRequest(cfg, cfg.VaultURL+cfg.VaultPath, "GET", nil, nil)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
			log.WithFields(log.Fields{
				"http_status":         resp.StatusCode,
				"http_status_message": resp.Status,
			}).Error("Invalid HTTP status received while fetching exsting keys")
			return fmt.Errorf("Invalid HTTP status received while fetching exsting keys")
		}
		if resp.StatusCode == http.StatusOK {
			err = json.Unmarshal(resp.Content, &res)
			if err != nil {
				return err
			}
			existing = res.Data
		}
	}

	existing[key] = value
	payload, err = json.Marshal(existing)
	if err != nil {
		return err
	}

	resp, err := httpRequest(cfg, cfg.VaultURL+cfg.VaultPath, "POST", nil, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	// Vault returns HTTP 204 on success
	if resp.StatusCode != http.StatusNoContent {
		log.WithFields(log.Fields{
			"http_status":         resp.StatusCode,
			"http_status_message": resp.Status,
		}).Error("Invalid HTTP status received")
		return fmt.Errorf("Setting of Vault key failed")
	}

	return nil
}
