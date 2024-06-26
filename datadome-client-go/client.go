package datadome

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// HOST_URL default datadome dashboard URL
const HOST_URL string = "https://customer-api.datadome.co/1.0/protection"

// Client to perform request on DataDome's API
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient is instantiate with given host and password parameters
func NewClient(host, password *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HOST_URL,
	}

	if host != nil {
		c.HostURL = *host
	}

	if password != nil {
		c.Token = *password
	}

	return &c, nil
}

// doRequest on the DataDome API with given http.Request and HttpResponse
func (c *Client) doRequest(req *http.Request, httpResponse *HttpResponse) (*HttpResponse, error) {
	// Add apikey as a query parameter on each request for authentication
	// Add also withoutTraffic parameter to true to have better performances
	q := req.URL.Query()
	q.Add("apikey", c.Token)
	q.Add("withoutTraffic", "true")
	req.URL.RawQuery = q.Encode()

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func () {
		err = res.Body.Close()
		if err != nil {
			log.Printf("failed to close body: %v", err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	log.Printf("[DEBUG] %s\n", body)

	err = json.Unmarshal(body, httpResponse)
	if err != nil {
		return nil, err
	}

	if httpResponse.Status != 200 {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, httpResponse.Errors)
	}

	return httpResponse, err
}
