package client

import (
	"encoding/json"
	"fmt"
)

// Role definitions
type Role struct {
	ID          int    `json:"roleId,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *Client) CreateRole(orgID string, role *Role) (*Role, error) {
	path := fmt.Sprintf("/org/%s/role", orgID)
	data, err := c.doRequest("PUT", path, role)
	if err != nil {
		return nil, err
	}
	var out Role
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) GetRole(orgID string, roleID int) (*Role, error) {
	path := fmt.Sprintf("/role/%d", roleID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var out Role
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) UpdateRole(orgID string, roleID int, role *Role) (*Role, error) {
	path := fmt.Sprintf("/role/%d", roleID)
	body := map[string]interface{}{
		"name":        role.Name,
		"description": role.Description,
	}
	data, err := c.doRequest("POST", path, body)
	if err != nil {
		return nil, err
	}
	var out Role
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) DeleteRole(orgID string, roleID int) error {
	path := fmt.Sprintf("/role/%d", roleID)
	// Workaround: Pangolin requires a replacement role ID for users in the deleted role.
	// We use ID 2 (Member) which is standard in a fresh org.
	body := map[string]interface{}{
		"roleId": "2",
	}
	_, err := c.doRequest("DELETE", path, body)
	return err
}

func (c *Client) ListRoles(orgID string) ([]Role, error) {
	path := fmt.Sprintf("/org/%s/roles", orgID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		Roles []Role `json:"roles"`
	}
	err = json.Unmarshal(data, &wrapper)
	return wrapper.Roles, err
}
