package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sort"
)

func cmdShow(cfg *Configuration, args []string) error {
	var res VaultKVResult
	var keys []string

	// remove duplicates
	temp := make(map[string]bool)
	for _, k := range args[1:] {
		temp[k] = true
	}
	for k := range temp {
		keys = append(keys, k)
	}
	sort.Strings(keys)

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

	if len(keys) == 0 {
		for k, v := range res.Data {
			fmt.Println(k, v)
		}
	} else if len(keys) == 1 {
		v, found := res.Data[keys[0]]
		if found {
			fmt.Println(v)
		}
	} else {
		for _, k := range keys {
			v, found := res.Data[k]
			if found {
				fmt.Println(k, v)
			}
		}
	}

	return nil
}
