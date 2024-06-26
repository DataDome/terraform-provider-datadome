package datadome

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// GetCustomRules list from the API
func (c *Client) GetCustomRules(ctx context.Context) ([]CustomRule, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/custom-rules", c.HostURL), nil)
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

	return customRules.CustomRules, nil
}


// CreateCutomRule with given CustomRule parameters
func (c *Client) CreateCustomRule(ctx context.Context, customRule CustomRule) (*ID, error) {
	reqBody := HttpRequest{
		Data: customRule,
	}
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/custom-rules", c.HostURL),
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

	return id, nil
}

// UpdateCustomRule by its ID
func (c *Client) UpdateCustomRule(ctx context.Context, customRule CustomRule) (*CustomRule, error) {
	reqBody := HttpRequest{
		Data: customRule,
	}
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %+v\n", customRule)

	req, err := http.NewRequestWithContext(
		ctx,
		"PUT",
		fmt.Sprintf("%s/custom-rules/%d", c.HostURL, customRule.ID),
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

	return &customRule, nil
}


// DeleteCustomRule by its ID
func (c *Client) DeleteCustomRule(ctx context.Context, customRule CustomRule) (*CustomRule, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"DELETE",
		fmt.Sprintf("%s/custom-rules/%d", c.HostURL, customRule.ID),
		nil,
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

	return &customRule, nil
}
