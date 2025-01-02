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

// HostURLCustomRule default datadome dashboard URL
const HostURLCustomRule string = "https://customer-api.datadome.co/1.1/protection/custom-rules"

// ClientCustomRule to perform request on DataDome's API
type ClientCustomRule struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClientCustomRules is instantiate with given host and password parameters
func NewClientCustomRules(host, password *string) (*ClientCustomRule, error) {
	c := ClientCustomRule{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURLCustomRule,
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
func (c *ClientCustomRule) doRequest(req *http.Request, httpResponse *HttpResponse) (*HttpResponse, error) {
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

// Read custom rule list from the API management
func (c *ClientCustomRule) Read(ctx context.Context, id int) (*CustomRule, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.HostURL, nil)
	if err != nil {
		return nil, err
	}

	customRules := &CustomRules{}
	resp := &HttpResponse{Data: customRules}

	_, err = c.doRequest(req, resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("response status is %d", resp.Status)
	}

	customRule := &CustomRule{}
	for _, v := range customRules.CustomRules {
		if v.ID == id {
			customRule = &v
		}
	}

	return customRule, nil
}

// Create custom rule with given CustomRule parameters
func (c *ClientCustomRule) Create(ctx context.Context, params CustomRule) (*int, error) {
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
	if resp.Status != 200 {
		return nil, fmt.Errorf("response status is %d", resp.Status)
	}

	return &id.ID, nil
}

// Update custom rule by its ID
func (c *ClientCustomRule) Update(ctx context.Context, params CustomRule) (*CustomRule, error) {
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
		"PUT",
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

// Delete custom rule by its ID
func (c *ClientCustomRule) Delete(ctx context.Context, id int) error {
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
	if resp.Status != 200 {
		return fmt.Errorf("response status is %d", resp.Status)
	}

	return nil
}
