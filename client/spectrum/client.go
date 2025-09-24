package spectrum

import (
	"crypto/tls"
	"net/http"
)

type Client struct {
	endpoint string
	username string
	password string

	httpClient *http.Client
}

func NewClient(endpoint string, us string, pw string, insecure bool) *Client {
	return &Client{
		endpoint: endpoint,
		username: us,
		password: pw,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecure,
				},
			},
		},
	}
}

func (c *Client) get(path string) ([]byte, error) {

	url := c.endpoint + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Username", c.username)
	req.Header.Add("X-Auth-Password", c.password)

	return nil, nil
}
