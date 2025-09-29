package unisphere

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

func NewClient(endpoint string, auth string, insecure bool) *UnisphereClient {
	return &UnisphereClient{
		endpoint: endpoint,
		auth:     auth,
		hc: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecure,
				},
			},
		},
	}
}

type StatusUnProcessableEntityError struct {
	Error struct {
		ErrorCode      int `json:"errorCode"`
		HttpStatusCode int `json:"httpStatusCode"`
		Messages       []struct {
			EnUS string `json:"en-US"`
		} `json:"messages"`
		Created time.Time `json:"created"`
	} `json:"error"`
}

func (c *UnisphereClient) get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.endpoint+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-EMC-REST-CLIENT", "true")
	// 최근 로그인 실패했을 경우, 인증
	if !c.lastAccess {
		if time.Since(c.accessTime).Minutes() < 1 {
			return nil, nil
		}
		c.hc.Jar, _ = cookiejar.New(nil)
		req.Header.Add("Authorization", "Basic "+c.auth)
	}
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check StatusCode
	var success bool
	switch resp.StatusCode {
	case http.StatusOK:
		success = true
	case http.StatusCreated:
		success = true
	case http.StatusAccepted:
		success = true
	case http.StatusUnauthorized:
		c.lastAccess = false
		err = errors.New("Unauthorized: " + string(body))
	case http.StatusUnprocessableEntity:
		var data StatusUnProcessableEntityError
		_ = json.Unmarshal(body, &data)
		err = errors.New("StatusUnprocessableEntity(422): " + data.Error.Messages[0].EnUS)
	}
	if !success {
		return nil, err
	}
	if !c.lastAccess {
		c.token = resp.Header.Get("EMC-CSRF-TOKEN")
		c.lastAccess = true
	}
	c.accessTime = time.Now()
	return body, nil

}

func (c *UnisphereClient) post(path string, data []byte) ([]byte, error) {
	if !c.lastAccess {
		return nil, nil
	}
	req, err := http.NewRequest("POST", c.endpoint+path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("EMC-CSRF-TOKEN", c.token)
	req.Header.Add("X-EMC-REST-CLIENT", "true")
	if !c.lastAccess {
		if time.Since(c.accessTime).Minutes() < 1 {
			return nil, nil
		}
		c.hc.Jar, _ = cookiejar.New(nil)
		req.Header.Add("Authorization", "Basic "+c.auth)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var success bool
	switch resp.StatusCode {
	case http.StatusOK:
		success = true
	case http.StatusCreated:
		success = true
	case http.StatusUnauthorized:
		c.lastAccess = success
		err = errors.New("Unauthorized: " + string(body))
	case http.StatusUnprocessableEntity:
		var respData StatusUnProcessableEntityError
		_ = json.Unmarshal(body, &respData)
		err = errors.New("StatusUnprocessableEntity(422): " + respData.Error.Messages[0].EnUS)
	}
	if !success {
		return nil, err
	}

	if !c.lastAccess {
		c.token = resp.Header.Get("EMC-CSRF-TOKEN")
		c.lastAccess = true
	}

	return body, nil
}
