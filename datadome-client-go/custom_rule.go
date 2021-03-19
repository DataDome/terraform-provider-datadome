package datadome

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetCustomRules() ([]CustomRule, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/custom-rules", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	customRules := resp.Data.CustomRules

	return customRules, nil
}

func (c *Client) CreateCustomRule(customRule CustomRule) (*CustomRule, error) {
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

	resp, err := c.doRequest(req)
	if err != nil || resp.Status != 200 {
		return nil, err
	}

	return &customRule, nil
}

func (c *Client) UpdateCustomRule(customRule CustomRule) (*CustomRule, error) {
	reqBody := HttpRequest{
		Data: customRule,
	}
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	customRulesList, err := c.GetCustomRules()
	if err != nil {
		return nil, err
	}

	customRuleId := -1
	for _, v := range customRulesList {
		if v.Name == customRule.Name {
			customRuleId = v.ID
		}
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/custom-rules/%d", c.HostURL, customRuleId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil || resp.Status != 200 {
		return nil, err
	}

	return &customRule, nil
}

func (c *Client) DeleteCustomRule(customRule CustomRule) (*CustomRule, error) {
	customRulesList, err := c.GetCustomRules()
	if err != nil {
		return nil, err
	}

	customRuleId := -1
	for _, v := range customRulesList {
		if v.Name == customRule.Name {
			customRuleId = v.ID
		}
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/custom-rules/%d", c.HostURL, customRuleId), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil || resp.Status != 200 {
		return nil, err
	}

	return &customRule, nil
}
