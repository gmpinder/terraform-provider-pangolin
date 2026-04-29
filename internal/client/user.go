package client

import (
	"encoding/json"
	"fmt"
)

// User definitions
type User struct {
	ID          string `json:"roleId,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *Client) CreateUser(orgID string, user *User) (*User, error) {
	path := fmt.Sprintf("/org/%s/user", orgID)
	data, err := c.doRequest("PUT", path, user)
	if err != nil {
		return nil, err
	}
	var out User
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) GetUser(orgID string, userID string) (*User, error) {
	path := fmt.Sprintf("/org/%s/user/%s", orgID, userID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var out User
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) DeleteUser(orgID string, userID int) error {
	path := fmt.Sprintf("/org/%s/user/%d", orgID, userID)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}
