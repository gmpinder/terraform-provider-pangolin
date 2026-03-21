package client

import (
	"encoding/json"
	"fmt"
)

type Resource struct {
	Name                  *string          `json:"name,omitempty"`
	Http                  *bool            `json:"http,omitempty"`
	Protocol              *string          `json:"protocol,omitempty"`
	EmailWhitelistEnabled *bool            `json:"emailWhitelistEnabled,omitempty"`
	Subdomain             *string          `json:"subdomain,omitempty"`
	ApplyRules            *bool            `json:"applyRules,omitempty"`
	DomainID              *string          `json:"domainId,omitempty"`
	ID                    *int64           `json:"resourceId,omitempty"`
	OrgID                 *string          `json:"orgId,omitempty"`
	NiceID                *string          `json:"niceId,omitempty"`
	Ssl                   *bool            `json:"ssl,omitempty"`
	BlockAccess           *bool            `json:"blockAccess,omitempty"`
	Sso                   *bool            `json:"sso,omitempty"`
	ProxyPort             *int32           `json:"proxyPort,omitempty"`
	Enabled               *bool            `json:"enabled,omitempty"`
	StickySession         *bool            `json:"stickySession,omitempty"`
	TlsServerName         *string          `json:"tlsServerName,omitempty"`
	SetHostHeader         *string          `json:"setHostHeader,omitempty"`
	Headers               []ResourceHeader `json:"headers,omitempty"`
	ProxyProtocol         *bool            `json:"proxyProtocol,omitempty"`
	ProxyProtocolVersion  *int32           `json:"proxyProtocolVersion,omitempty"`
	PostAuthPath          *string          `json:"postAuthPath,omitempty"`
}

type ResourceHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c *Client) CreateResource(orgID string, res *Resource) (*Resource, error) {
	path := fmt.Sprintf("/org/%s/resource", orgID)
	data, err := c.doRequest("PUT", path, res)
	if err != nil {
		return nil, err
	}
	var out Resource
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) GetResource(resID int64) (*Resource, error) {
	path := fmt.Sprintf("/resource/%d", resID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var out Resource
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) UpdateResource(resID int64, res *Resource) (*Resource, error) {
	path := fmt.Sprintf("/resource/%d", resID)
	data, err := c.doRequest("POST", path, res)
	if err != nil {
		return nil, err
	}
	var out Resource
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) DeleteResource(resID int64) error {
	path := fmt.Sprintf("/resource/%d", resID)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}
