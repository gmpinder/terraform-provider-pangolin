package client

import (
	"encoding/json"
	"fmt"
	"time"
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
	ID   int    `json:"siteId"`
	Name string `json:"name"`
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

func (c *Client) CreateSite(orgID string, name string) (*Site, error) {
	path := fmt.Sprintf("/org/%s/site", orgID)
	body := map[string]interface{}{
		"name":   name,
		"type":   "newt",
		"newtId": "test-newt-" + time.Now().Format("150405"),
		"secret": "test-secret-123",
	}
	data, err := c.doRequest("PUT", path, body)
	if err != nil {
		return nil, err
	}
	var out Site
	err = json.Unmarshal(data, &out)
	return &out, err
}
