package client

import (
	"encoding/json"
	"fmt"
)

// Organization definitions
type Organization struct {
	ID                              *string `json:"orgId,omitempty"`
	Name                            *string `json:"name,omitempty"`
	Subnet                          *string `json:"subnet,omitempty"`        // default: 100.90.128.0/24
	UtilitySubnet                   *string `json:"utilitySubnet,omitempty"` // default: 100.96.128.0/24
	RequireTwoFactor                *bool   `json:"requireTwoFactor,omitempty"`
	MaxSessionLengthHours           *int32  `json:"maxSessionLengthHours,omitempty"`
	PasswordExpiryDays              *int32  `json:"passwordExpiryDays,omitempty"`
	SettingsLogRetentionDaysRequest *int32  `json:"settingsLogRetentionDaysRequest,omitempty"`
	SettingsLogRetentionDaysAccess  *int32  `json:"settingsLogRetentionDaysAccess,omitempty"`
	SettingsLogRetentionDaysAction  *int32  `json:"settingsLogRetentionDaysAction,omitempty"`
}

func (c *Client) CreateOrganization(org *Organization) (*Organization, error) {
	data, err := c.doRequest("PUT", "/org", org)
	if err != nil {
		return nil, err
	}
	var out Organization
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) GetOrganization(orgID string) (*Organization, error) {
	path := fmt.Sprintf("/org/%s", orgID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var out Organization
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) UpdateOrganization(orgID string, org *Organization) (*Organization, error) {
	path := fmt.Sprintf("/org/%s", orgID)
	data, err := c.doRequest("POST", path, org)
	if err != nil {
		return nil, err
	}
	var out Organization
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) DeleteOrganization(orgID string) error {
	path := fmt.Sprintf("/org/%s", orgID)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}
