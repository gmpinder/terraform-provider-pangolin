package client

import (
	"encoding/json"
	"fmt"
)

// SiteResource definitions
type SiteResource struct {
	ID                 int      `json:"siteResourceId,omitempty"`
	NiceID             string   `json:"niceId,omitempty"`
	Name               string   `json:"name"`
	Mode               string   `json:"mode"`
	SiteID             int      `json:"siteId"`
	Destination        string   `json:"destination"`
	Enabled            bool     `json:"enabled"`
	Alias              *string  `json:"alias,omitempty"`
	UserIDs            []string `json:"userIds"`
	RoleIDs            []int    `json:"roleIds"`
	ClientIDs          []int    `json:"clientIds"`
	TCPPortRangeString string   `json:"tcpPortRangeString,omitempty"`
	UDPPortRangeString string   `json:"udpPortRangeString,omitempty"`
	DisableIcmp        bool     `json:"disableIcmp,omitempty"`
}

func (c *Client) CreateSiteResource(orgID string, res *SiteResource) (*SiteResource, error) {
	path := fmt.Sprintf("/org/%s/private-resource", orgID)
	body := map[string]interface{}{
		"name":        res.Name,
		"mode":        res.Mode,
		"siteId":      res.SiteID,
		"destination": res.Destination,
		"enabled":     res.Enabled,
		"userIds":     res.UserIDs,
		"roleIds":     res.RoleIDs,
		"clientIds":   res.ClientIDs,
	}
	if res.Alias != nil {
		body["alias"] = *res.Alias
	}
	data, err := c.doRequest("PUT", path, body)
	if err != nil {
		return nil, err
	}
	var out SiteResource
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) GetSiteResource(orgID string, siteID int, resID int) (*SiteResource, error) {
	path := fmt.Sprintf("/site-resource/%d", resID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var out SiteResource
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) UpdateSiteResource(resID int, res *SiteResource) (*SiteResource, error) {
	path := fmt.Sprintf("/site-resource/%d", resID)
	body := map[string]interface{}{
		"name":        res.Name,
		"siteId":      res.SiteID,
		"mode":        res.Mode,
		"destination": res.Destination,
		"enabled":     res.Enabled,
		"userIds":     res.UserIDs,
		"roleIds":     res.RoleIDs,
		"clientIds":   res.ClientIDs,
	}
	if res.Alias != nil {
		body["alias"] = *res.Alias
	}
	data, err := c.doRequest("POST", path, body)
	if err != nil {
		return nil, err
	}
	var out SiteResource
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) DeleteSiteResource(resID int) error {
	path := fmt.Sprintf("/site-resource/%d", resID)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}

func (c *Client) GetSiteResourceRoles(resID int) ([]int, error) {
	path := fmt.Sprintf("/site-resource/%d/roles", resID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Roles []struct {
			RoleID int `json:"roleId"`
		} `json:"roles"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}
	ids := make([]int, 0, len(wrapper.Roles))
	seen := make(map[int]bool)
	for _, r := range wrapper.Roles {
		if !seen[r.RoleID] {
			ids = append(ids, r.RoleID)
			seen[r.RoleID] = true
		}
	}
	return ids, nil
}

func (c *Client) GetSiteResourceUsers(resID int) ([]string, error) {
	path := fmt.Sprintf("/site-resource/%d/users", resID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Users []struct {
			UserID string `json:"userId"`
		} `json:"users"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}
	ids := make([]string, len(wrapper.Users))
	for i, u := range wrapper.Users {
		ids[i] = u.UserID
	}
	return ids, nil
}

func (c *Client) GetSiteResourceClients(resID int) ([]int, error) {
	path := fmt.Sprintf("/site-resource/%d/clients", resID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Clients []struct {
			ClientID int `json:"clientId"`
		} `json:"clients"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}
	ids := make([]int, len(wrapper.Clients))
	for i, cl := range wrapper.Clients {
		ids[i] = cl.ClientID
	}
	return ids, nil
}
