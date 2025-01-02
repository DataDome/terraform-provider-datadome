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
func (c *ClientEndpoint) doRequest(req *http.Request, httpResponse *HttpResponse) (*HttpResponse, error) {
	// Add apikey as a header on each request for authentication
	// Add also withoutTraffic parameter to true to have better performances
	q := req.URL.Query()
	req.Header.Set("x-api-key", c.Token)
	q.Add("withoutTraffic", "true")
	req.URL.RawQuery = q.Encode()

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
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

	if httpResponse.Status < 200 || httpResponse.Status > 299 {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, httpResponse.Errors)
	}

	return httpResponse, err
}

// Read endpoint information by its ID from the API management
func (c *ClientEndpoint) Read(ctx context.Context, id int) (*Endpoint, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	endpoint := &Endpoint{}
	resp := &HttpResponse{Data: endpoint}

	_, err = c.doRequest(req, resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("response status is %d", resp.Status)
	}

	return endpoint, nil
}

// Create new endpoint with given Endpoint parameters
func (c *ClientEndpoint) Create(ctx context.Context, params Endpoint) (*int, error) {
	reqBody := HttpRequest{
		Data: params,
	}
	rb, err := json.Marshal(reqBody)
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

	id := &ID{}
	resp := &HttpResponse{Data: id}

	resp, err = c.doRequest(req, resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != 201 {
		return nil, fmt.Errorf("response status is %d", resp.Status)
	}

	return &id.ID, nil
}

// Update endpoint by its ID
func (c *ClientEndpoint) Update(ctx context.Context, params Endpoint) (*Endpoint, error) {
	reqBody := HttpRequest{
		Data: params,
	}
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %+v\n", params)

	req, err := http.NewRequestWithContext(
		ctx,
		"PATCH",
		fmt.Sprintf("%s/%d", c.HostURL, params.ID),
		strings.NewReader(string(rb)),
	)
	if err != nil {
		return nil, err
	}

	resp := &HttpResponse{}

	resp, err = c.doRequest(req, resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("response status is %d", resp.Status)
	}

	return &params, nil
}

// Delete endpoint by its ID
func (c *ClientEndpoint) Delete(ctx context.Context, id int) error {
	req, err := http.NewRequestWithContext(
		ctx,
		"DELETE",
		fmt.Sprintf("%s/%d", c.HostURL, id),
		nil,
	)
	if err != nil {
		return err
	}

	resp := &HttpResponse{}

	resp, err = c.doRequest(req, resp)
	if err != nil {
		return err
	}
	if resp.Status != 204 {
		return fmt.Errorf("response status is %d", resp.Status)
	}

	return nil
}
