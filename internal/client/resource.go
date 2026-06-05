package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Resource struct {
	Name                  *string         `json:"name,omitempty"`
	Http                  *bool           `json:"http,omitempty"`
	Protocol              *string         `json:"protocol,omitempty"`
	EmailWhitelistEnabled *bool           `json:"emailWhitelistEnabled,omitempty"`
	Subdomain             *string         `json:"subdomain,omitempty"`
	ApplyRules            *bool           `json:"applyRules,omitempty"`
	DomainID              *string         `json:"domainId,omitempty"`
	ID                    *int64          `json:"resourceId,omitempty"`
	OrgID                 *string         `json:"orgId,omitempty"`
	NiceID                *string         `json:"niceId,omitempty"`
	Ssl                   *bool           `json:"ssl,omitempty"`
	BlockAccess           *bool           `json:"blockAccess,omitempty"`
	Sso                   *bool           `json:"sso,omitempty"`
	ProxyPort             *int32          `json:"proxyPort,omitempty"`
	Enabled               *bool           `json:"enabled,omitempty"`
	StickySession         *bool           `json:"stickySession,omitempty"`
	TlsServerName         *string         `json:"tlsServerName,omitempty"`
	SetHostHeader         *string         `json:"setHostHeader,omitempty"`
	Headers               ResourceHeaders `json:"headers,omitempty"`
	ProxyProtocol         *bool           `json:"proxyProtocol,omitempty"`
	ProxyProtocolVersion  *int32          `json:"proxyProtocolVersion,omitempty"`
	PostAuthPath          *string         `json:"postAuthPath,omitempty"`
}

type ResourceHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ResourceHeaders []ResourceHeader

type ResourceUsers struct {
	UserIds []string `json:"userIds"`
}

type ResourceRoles struct {
	RoleIds []int64 `json:"roleIds"`
}

const ADMIN_ID int64 = 1

// This is required because the `headers` property of the
// response during `PUT` is a stringified JSON array.
// This checks for that case and properly unmarshals it twice
// if it's a `string`.
func (th *ResourceHeaders) UnmarshalJSON(input []byte) error {
	if bytes.Equal(input, []byte("null")) {
		return nil
	}

	if len(input) > 1 {
		switch input[0] {
		case '"':
			var str string
			err := json.Unmarshal(input, &str)
			if err != nil {
				return err
			}

			err = readToResourceHeaders([]byte(str), th)
			if err != nil {
				return err
			}

			return nil
		case '[':
			err := readToResourceHeaders(input, th)
			if err != nil {
				return err
			}

			return nil
		default:
			break
		}
	}
	return fmt.Errorf("unable to read into ResourceHeader: %b", input)
}

func readToResourceHeaders(input []byte, th *ResourceHeaders) error {
	var headers []ResourceHeader
	err := json.Unmarshal(input, &headers)
	if err != nil {
		return err
	}

	*th = headers
	return nil
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

func (c *Client) UpdateResourceUsers(resID int64, users []string) error {
	path := fmt.Sprintf("/resource/%d/users", resID)
	payload := ResourceUsers{
		UserIds: users,
	}
	_, err := c.doRequest("POST", path, payload)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateResourceRoles(resID int64, roles []int64) error {
	path := fmt.Sprintf("/resource/%d/roles", resID)
	payload := ResourceRoles{
		RoleIds: roles,
	}
	_, err := c.doRequest("POST", path, payload)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetResourceUsers(resID int64) ([]string, error) {
	path := fmt.Sprintf("/resource/%d/users", resID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var out struct {
		Users []struct {
			UserId string `json:"userId"`
		} `json:"users"`
	}
	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}

	list := make([]string, len(out.Users))
	for i, value := range out.Users {
		list[i] = value.UserId
	}
	return list, nil
}

func (c *Client) GetResourceRoles(resID int64) ([]int64, error) {
	path := fmt.Sprintf("/resource/%d/roles", resID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var out struct {
		Roles []struct {
			RoleId int64 `json:"roleId"`
		} `json:"roles"`
	}
	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}

	list := make([]int64, 0, len(out.Roles))
	for _, value := range out.Roles {
		if value.RoleId != ADMIN_ID {
			list = append(list, value.RoleId)
		}
	}
	return list, nil
}

func (c *Client) DeleteResource(resID int64) error {
	path := fmt.Sprintf("/resource/%d", resID)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}

func (c *Client) GetEmailWhiteList(resdID int64) ([]string, error) {
	path := fmt.Sprintf("/resource/%d/whitelist", resdID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	type email struct {
		Email string `json:"email"`
	}
	var wrapper struct {
		Whitelist []email `json:"whitelist"`
	}
	err = json.Unmarshal(data, &wrapper)

	if err != nil {
		return nil, err
	}

	list := make([]string, len(wrapper.Whitelist))

	for i, value := range wrapper.Whitelist {
		list[i] = value.Email
	}
	return list, nil
}

func (c *Client) UpdateEmailWhitelist(resID int64, emailAddresses []string) error {
	path := fmt.Sprintf("/resource/%d/whitelist", resID)
	var wrapper struct {
		Emails []string `json:"emails"`
	}
	wrapper.Emails = emailAddresses
	_, err := c.doRequest("POST", path, wrapper)
	return err
}
