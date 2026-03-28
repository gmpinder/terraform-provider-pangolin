package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResource_BasicHTTP(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_resource.test", "name", "test-http-resource"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "http", "true"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "subdomain", "test-http"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "domain_id", "local"),
					resource.TestCheckResourceAttrSet("pangolin_resource.test", "id"),
					resource.TestCheckResourceAttrSet("pangolin_resource.test", "nice_id"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceConfigUpdated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_resource.test", "name", "test-http-resource-updated"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "ssl", "true"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "sticky_session", "true"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "block_access", "false"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "sso", "false"),
				),
			},
		},
	})
}

func TestAccResource_TCPNonHTTP(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfigTCPNonHTTP(8080),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_resource.test", "name", "test-tcp-resource"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "http", "false"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "proxy_port", "8080"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "proxy_protocol", "false"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceConfigTCPNonHTTPUpdated(9090),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_resource.test", "proxy_port", "9090"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "proxy_protocol", "true"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "proxy_protocol_version", "2"),
				),
			},
		},
	})
}

func TestAccResource_WithEmailWhitelist(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfigWithEmailWhitelist(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_resource.test", "email_whitelist_enabled", "true"),
					resource.TestCheckResourceAttr("pangolin_resource.test", "email_whitelist.#", "2"),
					resource.TestCheckTypeSetElemAttr("pangolin_resource.test", "email_whitelist.*", "test@example.com"),
					resource.TestCheckTypeSetElemAttr("pangolin_resource.test", "email_whitelist.*", "admin@example.com"),
				),
			},
			{
				Config: testAccResourceConfigWithEmailWhitelistUpdated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_resource.test", "email_whitelist.#", "1"),
					resource.TestCheckTypeSetElemAttr("pangolin_resource.test", "email_whitelist.*", "admin@example.com"),
				),
			},
		},
	})
}

func TestAccResource_WithHeaders(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfigWithHeaders(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_resource.test", "headers.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("pangolin_resource.test", "headers.*", map[string]string{
						"name":  "X-Custom-Header",
						"value": "custom-value",
					}),
				),
			},
		},
	})
}

func testAccResourceConfigBasic() string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_resource" "test" {
    org_id    = %[3]q
    name      = "test-http-resource"
    protocol  = "tcp"
    http      = true
    subdomain = "test-http"
    domain_id = "local"
}
`, testURL, testToken, testOrgID)
}

func testAccResourceConfigUpdated() string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_resource" "test" {
    org_id              = %[3]q
    name                = "test-http-resource-updated"
    protocol            = "tcp"
    http                = true
    subdomain           = "test-http-updated"
    domain_id           = "local"
    ssl                 = true
    sticky_session      = true
    block_access        = false
    sso                 = false
    enabled             = true
}
`, testURL, testToken, testOrgID)
}

func testAccResourceConfigTCPNonHTTP(port int) string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_resource" "test" {
    org_id            = %[3]q
    name              = "test-tcp-resource"
    protocol          = "tcp"
    http              = false
    proxy_port        = %[4]d
    enabled           = true
}
`, testURL, testToken, testOrgID, port)
}

func testAccResourceConfigTCPNonHTTPUpdated(port int) string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_resource" "test" {
    org_id                    = %[3]q
    name                      = "test-tcp-resource"
    protocol                  = "tcp"
    http                      = false
    proxy_port                = %[4]d
    proxy_protocol            = true
    proxy_protocol_version    = 2
    enabled                   = true
}
`, testURL, testToken, testOrgID, port)
}

func testAccResourceConfigWithEmailWhitelist() string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_resource" "test" {
    org_id                   = %[3]q
    name                     = "test-email-whitelist-resource"
    protocol                 = "tcp"
    http                     = true
    subdomain                = "test-email"
    domain_id                = "local"
    email_whitelist_enabled  = true
    email_whitelist          = ["test@example.com", "admin@example.com"]
    enabled                  = true
}
`, testURL, testToken, testOrgID)
}

func testAccResourceConfigWithEmailWhitelistUpdated() string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_resource" "test" {
    org_id                   = %[3]q
    name                     = "test-email-whitelist-resource"
    protocol                 = "tcp"
    http                     = true
    subdomain                = "test-email-updated"
    domain_id                = "local"
    email_whitelist_enabled  = true
    email_whitelist          = ["admin@example.com"]
    enabled                  = true
}
`, testURL, testToken, testOrgID)
}

func testAccResourceConfigWithHeaders() string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_resource" "test" {
    org_id            = %[3]q
    name              = "test-headers-resource"
    protocol          = "tcp"
    http              = true
    subdomain         = "test-headers"
    domain_id         = "local"
    enabled           = true

    headers = [
    	{
	        name  = "X-Custom-Header"
	        value = "custom-value"
	    },
		{
	        name  = "X-Another-Header"
	        value = "another-value"
	    }
	]
}
`, testURL, testToken, testOrgID)
}
