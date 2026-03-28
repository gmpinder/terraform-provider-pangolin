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

func (c *Client) DeleteResource(resID int64) error {
	path := fmt.Sprintf("/resource/%d", resID)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}
