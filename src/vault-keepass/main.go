package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var config Configuration
	var version = flag.Bool("version", false, "Show version")
	var help = flag.Bool("help", false, "Show help text")
	var path = flag.String("path", "", "Vault path")
	var timeout = flag.Int("timeout", 5, "Vault connection timeout in seconds")
	var insecure = flag.Bool("insecure-ssl", false, "Don't validate server certificate")
	var url = flag.String("url", "", "Vault URL")

	var _fmt = new(log.TextFormatter)
	_fmt.FullTimestamp = true
	_fmt.TimestampFormat = time.RFC3339
	log.SetFormatter(_fmt)

	flag.Usage = showUsage
	flag.Parse()

	if *version {
		showVersion()
		os.Exit(0)
	}

	if *help {
		showUsage()
		os.Exit(0)
	}

	if *path == "" {
		fmt.Fprintf(os.Stderr, "Error: Path to Vault key-value secrets engine is mandatory\n\n")
		showUsage()
		os.Exit(1)
	}
	config.VaultPath = filepath.Join("v1", *path)

	if *timeout <= 0 {
		fmt.Fprintf(os.Stderr, "Error: Timeout value must be greater than 0\n\n")
		os.Exit(1)
	}
	config.VaultTimeout = *timeout

	if *insecure {
		config.VaultInsecureSSL = true
	}

	setDefaults(&config)

	if *url != "" {
		config.VaultURL = *url
	}

	t, err := getVaultToken()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Unable to get current Vault token")
	}
	if t == "" {
		log.Fatal("Unable to get current Vault token")
	}
	config.VaultToken = t

	if config.VaultURL == "" {
		fmt.Fprintf(os.Stderr, "Error: No Vault URL found\n\n")
		showUsage()
		os.Exit(1)
	}

	err = checkConfiguration(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		os.Exit(1)
	}

	trailing := flag.Args()
	if len(trailing) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No command specified\n\n")
		showUsage()
		os.Exit(1)
	}

	switch trailing[0] {
	case "set":
		err := cmdSet(&config, trailing)
		if err != nil {
			log.Fatal(err)
		}

	case "list":
		err := cmdList(&config, trailing)
		if err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Fprintf(os.Stderr, "Error: Invalid command %s\n\n", trailing[0])
		showUsage()
		os.Exit(1)
	}

	os.Exit(0)
}
