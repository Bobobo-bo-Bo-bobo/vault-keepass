package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/url"
	"path/filepath"
)

func getVaultToken() (string, error) {
	// try environment variable VAULT_TOKEN first
	env := GetEnvironment()
	vlt, found := env["VAULT_TOKEN"]
	if found {
		return vlt, nil
	}

	// if HOME is set, try to read ${HOME}/
	home, found := env["HOME"]
	if found {
		content, err := ioutil.ReadFile(filepath.Join(home, ".vault-token"))
		if err != nil {
			log.WithFields(log.Fields{
				"file":  filepath.Join(home, ".vault-token"),
				"error": err.Error(),
			}).Error("Can't read token file")
			return "", err
		}
		return string(content), nil
	}

	return "", nil
}

func getDataFromVault(cfg *Configuration, u string) (*VaultKVResult, error) {
	var vkvrslt VaultKVResult

	result, err := httpRequest(cfg, u, "GET", nil, nil)
	if err != nil {
		return nil, err
	}

	switch result.StatusCode {
	case 200:
		break
	case 403:
		return nil, fmt.Errorf("Access denied (\"%s\") from Vault server. Are the token and/or the permissions to access %s valid?", result.Status, u)
	case 404:
		return nil, fmt.Errorf("Not found (\"%s\") from Vault server while accessing %s", result.Status, u)
	default:
		return nil, fmt.Errorf("Unexpected HTTP status, expected \"200 OK\" but got \"%s\" from %s instead", u, result.Status)
	}

	err = json.Unmarshal(result.Content, &vkvrslt)
	if err != nil {
		return nil, fmt.Errorf("%s (processing data from %s)", err.Error(), u)
	}

	return &vkvrslt, nil
}

func isVaultURL(s string) (bool, string, error) {
	if s == "" {
		return false, "", nil
	}

	parsed, err := url.Parse(s)
	if err != nil {
		return false, "", err
	}

	switch parsed.Scheme {
	case "http":
		if len(parsed.Path) > 0 {
			if parsed.Path[len(parsed.Path)-1] == '/' {
				parsed.Path = parsed.Path[0 : len(parsed.Path)-1]
			}
		}
		return true, fmt.Sprintf("http://%s/%s", parsed.Host, parsed.Path), nil

	case "https":
		if len(parsed.Path) > 0 {
			if parsed.Path[len(parsed.Path)-1] == '/' {
				parsed.Path = parsed.Path[0 : len(parsed.Path)-1]
			}
		}
		return true, fmt.Sprintf("https://%s/%s", parsed.Host, parsed.Path), nil

	default:
		return false, "", fmt.Errorf("Unrecognized URL scheme %s", parsed.Scheme)
	}

}
