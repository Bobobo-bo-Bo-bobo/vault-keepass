package main

import "fmt"

const name = "vault-keepass"
const version = "1.0.0-2020.01.25"
const _url = "https://git.ypbind.de/cgit/vault-keepass/"

var userAgent = fmt.Sprintf("%s/%s (%s)", name, version, _url)

const helpText = `
Usage: vault-keepass [-help] [-url <url>] [-version] -path <path> <command>
    -help           This text

    -insecure-ssl   Don't validate server certificate

    -path <path>    Path to Vault key value backend

    -timeout <sec>  Connection timeout in seconds.
                    Default: 5

    -url <url>      Use <url> as vault URL instead of value from
                    environment variable VAULT_ADDR or the default
                    value of http://localhost:8200

    -version        Show version information

Commands:
    copy            Copy value to clipboard

    delete          Delete a key/value entry

    list            List keys

    show            Display value on standard output

    set [-replace] key value
        Set a key to a value

            -replace    Replace whole content on path with key=value
`
