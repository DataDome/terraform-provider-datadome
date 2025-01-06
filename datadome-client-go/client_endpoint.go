package datadome

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// HOST_URL default datadome dashboard URL
const HostURLEndpoint string = "https://customer-api.datadome.co/1.0/endpoints"

// ClientEndpoint to perform request on DataDome's API
type ClientEndpoint struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClientEndpoint is instantiate with given host and password parameters
func NewClientEndpoint(host, password *string) (*ClientEndpoint, error) {
	c := ClientEndpoint{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURLEndpoint,
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
func (c *ClientEndpoint) doRequest(req *http.Request, endpoint *Endpoint) error {
	// Add apikey as a header on each request for authentication
	req.Header.Set("x-api-key", c.Token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			log.Printf("failed to close body: %v", err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	log.Printf("[DEBUG] %s\n", body)

	if endpoint != nil {
		err = json.Unmarshal(body, endpoint)
		if err != nil {
			return err
		}
	}

	return err
}

// Read endpoint information by its ID from the API management
func (c *ClientEndpoint) Read(ctx context.Context, id string) (*Endpoint, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	endpoint := &Endpoint{}

	err = c.doRequest(req, endpoint)
	if err != nil {
		return nil, err
	}

	return endpoint, nil
}

// Create new endpoint with given Endpoint parameters
func (c *ClientEndpoint) Create(ctx context.Context, params Endpoint) (*string, error) {
	rb, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.HostURL,
		strings.NewReader(string(rb)),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	endpoint := &Endpoint{}

	err = c.doRequest(req, endpoint)
	if err != nil {
		return nil, err
	}

	return endpoint.ID, nil
}

// Update endpoint by its ID
func (c *ClientEndpoint) Update(ctx context.Context, params Endpoint) (*Endpoint, error) {
	rb, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %+v\n", params)

	req, err := http.NewRequestWithContext(
		ctx,
		"PATCH",
		fmt.Sprintf("%s/%s", c.HostURL, *params.ID),
		strings.NewReader(string(rb)),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/merge-patch+json")

	endpoint := &Endpoint{}

	err = c.doRequest(req, endpoint)
	if err != nil {
		return nil, err
	}

	return endpoint, nil
}

// Delete endpoint by its ID
func (c *ClientEndpoint) Delete(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(
		ctx,
		"DELETE",
		fmt.Sprintf("%s/%s", c.HostURL, id),
		nil,
	)
	if err != nil {
		return err
	}

	err = c.doRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
}
