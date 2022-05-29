package vault

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type SecretGetter interface {
	Get(path string) (*string, error)
}

type Client struct {
	Url string
	Token string

	httpClient *http.Client
}

func NewClient(url, token string) Client {
	return Client{
		Url: strings.TrimRight(url, "/"),
		Token: token,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (c Client) Get(path string) (*string, error) {
	req, err := c.getVaultRequest("GET", c.trimPath(path), nil)
	if err != nil {
		return nil, err
	}

	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// read out the body of the request
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	// decode the body
	bodyString := string(body)

	// catch a non-successful http code
	if rsp.StatusCode != 200 {
		return &bodyString, errors.New(fmt.Sprint("Vault Client: Error in GET to `", path, "`. Response Code: ", rsp.StatusCode))
	} else {
		// full on success, return nil for error
		return &bodyString, nil
	}
}

func (c Client) trimPath(path string) string {
	return strings.TrimLeft(path, "/")
}

func (c Client) getVaultRequest(method, path string, body io.Reader) (*http.Request, error) {
	// create a new http request with the vault token already provided in the headers
	if req, err := http.NewRequest(method, strings.Join([]string{c.Url, path}, "/"), body); err == nil {
		req.Header.Add("X-Vault-Token", c.Token)
		return req, nil
	} else {
		return nil, err
	}
}
