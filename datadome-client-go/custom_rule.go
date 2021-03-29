package datadome

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (c *Client) GetCustomRules() ([]CustomRule, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/custom-rules", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	customRules := &CustomRules{}
	resp := &HttpResponse{Data: customRules}

	resp, err = c.doRequest(req, resp)

	if err != nil {
		return nil, err
	}

	return customRules.CustomRules, nil
}

func (c *Client) CreateCustomRule(customRule CustomRule) (*ID, error) {
	reqBody := HttpRequest{
		Data: customRule,
	}
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/custom-rules", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	id := &ID{}
	resp := &HttpResponse{Data: id}

	resp, err = c.doRequest(req, resp)
	if err != nil || resp.Status != 200 {
		return nil, err
	}

	return id, nil
}

func (c *Client) UpdateCustomRule(customRule CustomRule) (*CustomRule, error) {
	reqBody := HttpRequest{
		Data: customRule,
	}
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %+v\n", customRule)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/custom-rules/%d", c.HostURL, customRule.ID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	resp := &HttpResponse{}

	resp, err = c.doRequest(req, resp)
	if err != nil || resp.Status != 200 {
		return nil, err
	}

	return &customRule, nil
}

func (c *Client) DeleteCustomRule(customRule CustomRule) (*CustomRule, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/custom-rules/%d", c.HostURL, customRule.ID), nil)
	if err != nil {
		return nil, err
	}

	resp := &HttpResponse{}

	resp, err = c.doRequest(req, resp)
	if err != nil || resp.Status != 200 {
		return nil, err
	}

	return &customRule, nil
}
