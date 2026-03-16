package client

import (
	"encoding/json"
	"fmt"
)

// Organization definitions
type Organization struct {
	OrgId                           string  `json:"orgId"`
	Name                            string  `json:"name"`
	Subnet                          string  `json:"subnet"`
	UtilitySubnet                   string  `json:"utilitySubnet"`
	CreatedAt                       *string `json:"createdAt,omitempty"`
	RequireTwoFactor                *bool   `json:"requireTwoFactor,omitempty"`
	MaxSessionLengthHours           *int    `json:"maxSessionLengthHours,omitempty"`
	PasswordExpiryDays              *int    `json:"passwordExpiryDays,omitempty"`
	SettingsLogRetentionDaysRequest *int    `json:"settingsLogRetentionDaysRequest,omitempty"`
	SettingsLogRetentionDaysAccess  *int    `json:"settingsLogRetentionDaysAccess,omitempty"`
	SettingsLogRetentionDaysAction  *int    `json:"settingsLogRetentionDaysAction,omitempty"`
	SshCaPrivateKey                 *string `json:"sshCaPrivateKey,omitempty"`
	SshCaPublicKey                  *string `json:"sshCaPublicKey,omitempty"`
	IsBillingOrg                    *bool   `json:"isBillingOrg,omitempty"`
	BillingOrgId                    *string `json:"billingOrgId,omitempty"`
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

func (c *Client) DeleteOrganization(orgID string) error {
	path := fmt.Sprintf("/org/%s", orgID)
	_, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ListOrganizations() ([]Organization, error) {
	data, err := c.doRequest("GET", "/orgs?limit=1000&offset=0", nil)
	if err != nil {
		return nil, err
	}

	var out struct {
		Orgs []Organization `json:"orgs"`
	}
	err = json.Unmarshal(data, &out)
	return out.Orgs, err
}
