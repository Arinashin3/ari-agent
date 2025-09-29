package spectrum

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Client struct {
	endpoint   string
	username   string
	password   string
	token      string
	lastAuth   bool
	lastAccess time.Time

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

func (c *Client) login() error {
	if c.lastAuth {
		return nil
	}
	if time.Since(c.lastAccess) < time.Minute {
		return nil
	}
	url := c.endpoint + "/rest/auth"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Username", c.username)
	req.Header.Add("X-Auth-Password", c.password)
	c.httpClient.Jar, _ = cookiejar.New(nil)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("Invalid Status Code(" + resp.Status + ") :" + string(body))
	}
	if err != nil {
		return err
	}
	var loginResp struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return err
	}
	c.token = loginResp.Token

	c.lastAuth = true
	c.lastAccess = time.Now()
	return nil
}

func (c *Client) post(path string, reqBody []byte) ([]byte, error) {
	err := c.login()
	if err != nil {
		return nil, err
	}
	url := c.endpoint + path
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Invalid status code: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	c.lastAccess = time.Now()
	return body, nil
}
