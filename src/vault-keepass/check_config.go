package main

import "fmt"

func checkConfiguration(cfg *Configuration) error {
	is, u, err := isVaultURL(cfg.VaultURL)
	if err != nil {
		return err
	}
	if !is {
		return fmt.Errorf("Not a valid Vault URL")
	}
	cfg.VaultURL = u
	return nil
}
