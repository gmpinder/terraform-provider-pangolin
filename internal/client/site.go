package client

import (
	"encoding/json"
	"fmt"
)

// Site defaults
type SiteDefaults struct {
	ExitNodeId    int    `json:"exitNodeId"`
	Address       string `json:"address"`
	PublicKey     string `json:"publicKey"`
	Name          string `json:"name"`
	ListenPort    int    `json:"listenPort"`
	Endpoint      string `json:"endpoint"`
	Subnet        string `json:"subnet"`
	ClientAddress string `json:"clientAddress"`
	NewtId        string `json:"newtId"`
	NewtSecret    string `json:"newtSecret"`
}

func (c *Client) GetSiteDefaults(orgID string) (*SiteDefaults, error) {
	path := fmt.Sprintf("/org/%s/pick-site-defaults", orgID)
	data, err := c.doRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}
	var sd SiteDefaults
	err = json.Unmarshal(data, &sd)
	return &sd, err
}

// Site definitions
type Site struct {
	ID      *int64  `json:"siteId,omitempty"`
	OrgID   *string `json:"orgId,omitempty"`
	Name    *string `json:"name,omitempty"`
	NiceId  *string `json:"niceId,omitempty"`
	NewtID  *string `json:"newtId,omitempty"`
	PubKey  *string `json:"pubKey,omitempty"`
	Secret  *string `json:"secret,omitempty"`
	Address *string `json:"address,omitempty"`
	Subnet  *string `json:"subnet,omitempty"`
	Type    *string `json:"type,omitempty"`
}

func (c *Client) ListSites(orgID string) ([]Site, error) {
	path := fmt.Sprintf("/org/%s/sites", orgID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Sites []Site `json:"sites"`
	}
	err = json.Unmarshal(data, &wrapper)
	return wrapper.Sites, err
}

func (c *Client) GetSite(orgID string, siteID int64) (*Site, error) {
	path := fmt.Sprintf("/site/%d", siteID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var out Site
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) UpdateSite(siteID int64, site *Site) (*Site, error) {
	path := fmt.Sprintf("/site/%d", siteID)
	data, err := c.doRequest("POST", path, site)
	if err != nil {
		return nil, err
	}
	var out Site
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) CreateSite(orgID string, site *Site) (*Site, error) {
	path := fmt.Sprintf("/org/%s/site", orgID)
	data, err := c.doRequest("PUT", path, site)
	if err != nil {
		return nil, err
	}
	var out Site
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) DeleteSite(siteID int64) error {
	path := fmt.Sprintf("/site/%d", siteID)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}
