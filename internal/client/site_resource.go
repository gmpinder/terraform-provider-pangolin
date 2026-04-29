package client

import (
	"encoding/json"
	"fmt"
)

// SiteResource definitions
type SiteResource struct {
	ID                 *int64   `json:"siteResourceId,omitempty"`
	NiceID             *string  `json:"niceId,omitempty"`
	OrgID              *string  `json:"orgId,omitempty"`
	Name               *string  `json:"name,omitempty"`
	Mode               *string  `json:"mode,omitempty"`
	SiteIDs            []int64  `json:"siteIds"`
	Destination        *string  `json:"destination,omitempty"`
	Enabled            *bool    `json:"enabled,omitempty"`
	Alias              *string  `json:"alias,omitempty"`
	UserIDs            []string `json:"userIds"`
	RoleIDs            []int64  `json:"roleIds"`
	ClientIDs          []int64  `json:"clientIds"`
	TCPPortRangeString *string  `json:"tcpPortRangeString,omitempty"`
	UDPPortRangeString *string  `json:"udpPortRangeString,omitempty"`
	DisableIcmp        *bool    `json:"disableIcmp,omitempty"`
}

func (c *Client) CreateSiteResource(orgID string, res *SiteResource) (*SiteResource, error) {
	path := fmt.Sprintf("/org/%s/site-resource", orgID)
	data, err := c.doRequest("PUT", path, res)
	if err != nil {
		return nil, err
	}
	var out SiteResource
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) ListSiteResources(orgID string) ([]SiteResource, error) {
	path := fmt.Sprintf("/org/%s/site-resources", orgID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		SiteResources []SiteResource `json:"siteResources"`
	}
	err = json.Unmarshal(data, &out)
	return out.SiteResources, err
}

func (c *Client) GetSiteResource(orgID string, resID int64) (*SiteResource, error) {
	srList, err := c.ListSiteResources(orgID)
	if err != nil {
		return nil, err
	}

	for _, sr := range srList {
		if *sr.ID == resID {
			return &sr, nil
		}
	}
	return nil, fmt.Errorf("failed to find site resource %d for org %s", resID, orgID)
}

func (c *Client) UpdateSiteResource(resID int64, res *SiteResource) (*SiteResource, error) {
	path := fmt.Sprintf("/site-resource/%d", resID)
	data, err := c.doRequest("POST", path, res)
	if err != nil {
		return nil, err
	}
	var out SiteResource
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) DeleteSiteResource(resID int64) error {
	path := fmt.Sprintf("/site-resource/%d", resID)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}

func (c *Client) GetSiteResourceRoles(resID int64) ([]int, error) {
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

func (c *Client) GetSiteResourceUsers(resID int64) ([]string, error) {
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

func (c *Client) GetSiteResourceClients(resID int64) ([]int, error) {
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
