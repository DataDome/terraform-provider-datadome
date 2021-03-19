package datadome

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Default datadome dashboard URL
const HostURL string = "https://dev-app.datadome.co/1.0/protection"

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

func NewClient(host, password *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	if password != nil {
		c.Token = *password
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) (*HttpResponse, error) {
	// Add apikey as a query parameter on each request for authentication
	q := req.URL.Query()
	q.Add("apikey", c.Token)
	req.URL.RawQuery = q.Encode()

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	httpResponse := HttpResponse{}
	err = json.Unmarshal(body, &httpResponse)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK || httpResponse.Status != 200 {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, httpResponse.Errors)
	}

	return &httpResponse, err
}
