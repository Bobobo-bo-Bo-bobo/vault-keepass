package main

import (
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func httpRequest(config *Configuration, _url string, method string, header *map[string]string, reader io.Reader) (HTTPResult, error) {
	var result HTTPResult
	var transp *http.Transport

	if config.VaultInsecureSSL {
		transp = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		transp = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
	}

	client := &http.Client{
		Timeout:   time.Duration(config.VaultTimeout) * time.Second,
		Transport: transp,
		// non-GET methods (like PATCH, POST, ...) may or may not work when encountering
		// HTTP redirect. Don't follow 301/302. The new location can be checked by looking
		// at the "Location" header.
		CheckRedirect: func(http_request *http.Request, http_via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	request, err := http.NewRequest(method, _url, reader)
	if err != nil {
		return result, err
	}

	defer func() {
		if request.Body != nil {
			ioutil.ReadAll(request.Body)
			request.Body.Close()
		}
	}()

	// set User-Agent
	request.Header.Set("User-Agent", userAgent)

	// X-Vault-Token is mandatory
	if config.VaultToken == "" {
		return result, fmt.Errorf("No token for Vault access")
	}
	request.Header.Set("X-Vault-Token", config.VaultToken)
	request.Header.Set("X-Vault-Request", "true")

	// close connection after response and prevent re-use of TCP connection because some implementations (e.g. HP iLO4)
	// don't like connection reuse and respond with EoF for the next connections
	request.Close = true

	// add supplied additional headers
	if header != nil {
		for key, value := range *header {
			request.Header.Add(key, value)
		}
	}

	if config.Debug {
		log.WithFields(log.Fields{
			"url":     _url,
			"method":  method,
			"headers": request.Header,
		}).Debug("Sending HTTP request")

	}

	response, err := client.Do(request)
	if err != nil {
		return result, err
	}

	defer func() {
		ioutil.ReadAll(response.Body)
		response.Body.Close()
	}()

	result.Status = response.Status
	result.StatusCode = response.StatusCode
	result.Header = response.Header
	result.Content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	return result, nil
}
