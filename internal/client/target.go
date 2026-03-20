package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Target definitions
type Target struct {
	ID                  *int64        `json:"targetId,omitempty"`
	SiteID              *int64        `json:"siteId,omitempty"`
	ResourceID          *int64        `json:"resourceId,omitempty"`
	IP                  *string       `json:"ip,omitempty"`
	Port                *int32        `json:"port,omitempty"`
	Method              *string       `json:"method,omitempty"`
	Enabled             *bool         `json:"enabled,omitempty"`
	HCEnabled           *bool         `json:"hcEnabled,omitempty"`
	HCPath              *string       `json:"hcPath,omitempty"`
	HCScheme            *string       `json:"hcScheme,omitempty"`
	HCMode              *string       `json:"hcMode,omitempty"`
	HCHostname          *string       `json:"hcHostname,omitempty"`
	HCPort              *int32        `json:"hcPort,omitempty"`
	HCInterval          *int32        `json:"hcInterval,omitempty"`
	HCUnhealthyInterval *int32        `json:"hcUnhealthyInterval,omitempty"`
	HCTimeout           *int32        `json:"hcTimeout,omitempty"`
	HCHeaders           TargetHeaders `json:"hcHeaders,omitempty"`
	HCFollowRedirects   *bool         `json:"hcFollowRedirects,omitempty"`
	HCMethod            *string       `json:"hcMethod,omitempty"`
	HCStatus            *int32        `json:"hcStatus,omitempty"`
	HCTlsServerName     *string       `json:"hcTlsServerName,omitempty"`
	Path                *string       `json:"path,omitempty"`
	PathMatchType       *string       `json:"pathMatchType,omitempty"`
	RewritePath         *string       `json:"rewritePath,omitempty"`
	RewritePathType     *string       `json:"rewritePathType,omitempty"`
	Priority            *int32        `json:"priority,omitempty"`
}

type TargetHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TargetHeaders []TargetHeader

// This is required because the `hcHeaders` property of the
// response during `PUT` is a stringified JSON array.
// This checks for that case and properly unmarshals it twice
// if it's a `string`.
func (th *TargetHeaders) UnmarshalJSON(input []byte) error {
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

			err = readToTargetHeaders([]byte(str), th)
			if err != nil {
				return err
			}

			return nil
		case '[':
			err := readToTargetHeaders(input, th)
			if err != nil {
				return err
			}

			return nil
		default:
			break
		}
	}
	return fmt.Errorf("unable to read into TargetHeaders: %b", input)
}

func readToTargetHeaders(input []byte, th *TargetHeaders) error {
	var headers []TargetHeader
	err := json.Unmarshal(input, &headers)
	if err != nil {
		return err
	}

	*th = headers
	return nil
}

func (c *Client) CreateTarget(resID int64, target *Target) (*Target, error) {
	path := fmt.Sprintf("/resource/%d/target", resID)
	// Add other optional fields if needed...
	data, err := c.doRequest("PUT", path, target)
	if err != nil {
		return nil, err
	}
	var out Target
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) GetTarget(targetID int64) (*Target, error) {
	path := fmt.Sprintf("/target/%d", targetID)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var out Target
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) UpdateTarget(targetID int64, target *Target) (*Target, error) {
	path := fmt.Sprintf("/target/%d", targetID)
	data, err := c.doRequest("POST", path, target)
	if err != nil {
		return nil, err
	}
	var out Target
	err = json.Unmarshal(data, &out)
	return &out, err
}

func (c *Client) DeleteTarget(targetID int64) error {
	path := fmt.Sprintf("/target/%d", targetID)
	_, err := c.doRequest("DELETE", path, nil)
	return err
}
