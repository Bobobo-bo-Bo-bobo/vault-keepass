package main

import "fmt"

const name = "vault-keepass"
const version = "1.0.0-2020.02.24"
const _url = "https://git.ypbind.de/cgit/vault-keepass/"

var userAgent = fmt.Sprintf("%s/%s (%s)", name, version, _url)

const helpText = `
Usage: vault-keepass [-debug] [-help] [-url <url>] [-version] -path <path> <command>
    -debug          Enable debug mode.
                    WARNING: This will expose the vault token and all secrets!

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
    copy <key>
        Copy value of <key> to clipboard.
        Only supported on Linux, Windows and MacOS

    delete <loc> [<loc> ...]
        Delete all key entries from Vault storage at <loc> in the provided path

    list
        List keys in the path provided

    show [<key>] [<key>] ...
        Display key value on standard output
        List of keys to display, if empty all keys and values will be displayed.
        If only one key is given it's value (without key name) will be displayed.

    set [-replace] [key value]
        Set a key to a value

            -replace    Replace whole content on path with key=value

        If key and value are not provided they will be read from standard input.
`
