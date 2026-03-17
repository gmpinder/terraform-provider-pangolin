package client

import (
	"encoding/json"
	"fmt"
)

// Rule definitions
type Rule struct {
	ID         *int64 `json:"ruleId,omitempty"`
	ResourceID *int64 `json:"resourceId,omitempty"`
	Action     string `json:"action"`
	Match      string `json:"match"`
	Value      string `json:"value"`
	Priority   int32  `json:"priority"`
	Enabled    bool   `json:"enabled"`
}

func (c *Client) CreateRule(resourceID int64, rule *Rule) (*Rule, error) {
	path := fmt.Sprintf("/resource/%d/rule", resourceID)
	data, err := c.doRequest("PUT", path, rule)
	if err != nil {
		return nil, err
	}
	var out Rule
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) UpdateRule(resourceID int64, ruleID int64, rule *Rule) (*Rule, error) {
	path := fmt.Sprintf("/resource/%d/rule/%d", resourceID, ruleID)
	data, err := c.doRequest("POST", path, rule)
	if err != nil {
		return nil, err
	}
	var out Rule
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) ListRules(resourceID int64) ([]Rule, error) {
	path := fmt.Sprintf("/resource/%d/rules", resourceID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Rules []Rule `json:"rules"`
	}
	err = json.Unmarshal(data, &out)
	return out.Rules, err
}

func (c *Client) GetRule(resourceID int64, ruleID int64) (*Rule, error) {
	rules, err := c.ListRules(resourceID)
	if err != nil {
		return nil, err
	}

	for _, rule := range rules {
		if *rule.ID == ruleID {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("failed to get rule %d on resource %d", ruleID, resourceID)
}

func (c *Client) DeleteRule(resourceID int64, ruleID int64) error {
	path := fmt.Sprintf("/resource/%d/rule/%d", resourceID, ruleID)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}
