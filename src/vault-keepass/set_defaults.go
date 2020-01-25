package main

func setDefaults(cfg *Configuration) {
	env := GetEnvironment()
	u, f := env["VAULT_ADDR"]
	if f {
		cfg.VaultURL = u
	}
}
