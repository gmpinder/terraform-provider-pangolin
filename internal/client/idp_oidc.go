package client

import (
	"encoding/json"
	"fmt"
)

// OIDCIdP represents an OIDC Identity Provider
type OIDCIdP struct {
	ID             *int64  `json:"idpId,omitempty"`
	Name           *string `json:"name,omitempty"`
	ClientID       *string `json:"clientId,omitempty"`
	ClientSecret   *string `json:"clientSecret,omitempty"`
	AuthURL        *string `json:"authUrl,omitempty"`
	TokenURL       *string `json:"tokenUrl,omitempty"`
	IdentifierPath *string `json:"identifierPath,omitempty"`
	EmailPath      *string `json:"emailPath,omitempty"`
	NamePath       *string `json:"namePath,omitempty"`
	Scopes         *string `json:"scopes,omitempty"`
	AutoProvision  *bool   `json:"autoProvision,omitempty"`
	Tags           *string `json:"tags,omitempty"`
	Variant        *string `json:"variant,omitempty"`
}

// CreateOIDCIdP creates a new OIDC IdP
func (c *Client) CreateOIDCIdP(idp *OIDCIdP) (*OIDCIdP, error) {
	path := "/idp/oidc"
	data, err := c.doRequest("PUT", path, idp)
	if err != nil {
		return nil, err
	}
	var out OIDCIdP
	err = json.Unmarshal(data, &out)
	return &out, err
}

// GetOIDCIdP retrieves an OIDC IdP by ID
func (c *Client) GetOIDCIdP(idpID int64) (*OIDCIdP, error) {
	path := fmt.Sprintf("/idp/%d", idpID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var out OIDCIdP
	err = json.Unmarshal(data, &out)
	return &out, err
}

// UpdateOIDCIdP updates an existing OIDC IdP
func (c *Client) UpdateOIDCIdP(idpID int64, idp *OIDCIdP) (*OIDCIdP, error) {
	path := fmt.Sprintf("/idp/%d/oidc", idpID)
	data, err := c.doRequest("POST", path, idp)
	if err != nil {
		return nil, err
	}
	var out OIDCIdP
	err = json.Unmarshal(data, &out)
	return &out, err
}

// DeleteOIDCIdP deletes an OIDC IdP by ID
func (c *Client) DeleteOIDCIdP(idpID int64) error {
	path := fmt.Sprintf("/idp/%d", idpID)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}
