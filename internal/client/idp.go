package client

import (
	"encoding/json"
	"fmt"
)

// Idp represents an OIDC Identity Provider
type Idp struct {
	ID                 *int64  `json:"idpId,omitempty"`
	Name               *string `json:"name,omitempty"`
	AutoProvision      *bool   `json:"autoProvision,omitempty"`
	RedirectUrl        *string `json:"redirectUrl,omitempty"`
	DefaultRoleMapping *string `json:"defaultRoleMapping,omitempty"`
	DefaultOrgMapping  *string `json:"defaultOrgMapping,omitempty"`
}

type IdpOidcConfig struct {
	ClientID       *string `json:"clientId,omitempty"`
	ClientSecret   *string `json:"clientSecret,omitempty"`
	AuthURL        *string `json:"authUrl,omitempty"`
	TokenURL       *string `json:"tokenUrl,omitempty"`
	IdentifierPath *string `json:"identifierPath,omitempty"`
	EmailPath      *string `json:"emailPath,omitempty"`
	NamePath       *string `json:"namePath,omitempty"`
	Scopes         *string `json:"scopes,omitempty"`
	Variant        *string `json:"variant,omitempty"`
}

type IdpPost struct {
	Idp
	IdpOidcConfig
}

type IdpGet struct {
	Idp           *Idp           `json:"idp,omitempty"`
	IdpOidcConfig *IdpOidcConfig `json:"idpOidcConfig,omitempty"`
}

// CreateIdP creates a new OIDC IdP
func (c *Client) CreateIdP(idp *IdpPost) (*Idp, error) {
	path := "/idp/oidc"
	data, err := c.doRequest("PUT", path, idp)
	if err != nil {
		return nil, err
	}
	var out Idp
	err = json.Unmarshal(data, &out)
	return &out, err
}

// GetIdP retrieves an OIDC IdP by ID
func (c *Client) GetIdP(idpID int64) (*IdpGet, error) {
	path := fmt.Sprintf("/idp/%d", idpID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var out IdpGet
	err = json.Unmarshal(data, &out)
	return &out, err
}

// UpdateIdP updates an existing OIDC IdP
func (c *Client) UpdateIdP(idpID int64, idp *IdpPost) (*Idp, error) {
	path := fmt.Sprintf("/idp/%d/oidc", idpID)
	data, err := c.doRequest("POST", path, idp)
	if err != nil {
		return nil, err
	}
	var out Idp
	err = json.Unmarshal(data, &out)
	return &out, err
}

// DeleteOIDCIdP deletes an OIDC IdP by ID
func (c *Client) DeleteOIDCIdP(idpID int64) error {
	path := fmt.Sprintf("/idp/%d", idpID)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}
