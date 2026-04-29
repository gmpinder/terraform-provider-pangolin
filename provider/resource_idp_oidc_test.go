package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIDPOidc_Basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIDPOidcConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_idp_oidc.test", "name", "test-idp"),
					resource.TestCheckResourceAttr("pangolin_idp_oidc.test", "client_id", "test-client-id"),
					resource.TestCheckResourceAttr("pangolin_idp_oidc.test", "client_secret", "test-client-secret"),
					resource.TestCheckResourceAttrSet("pangolin_idp_oidc.test", "id"),
				),
			},
			{
				Config: testAccIDPOidcConfigUpdated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_idp_oidc.test", "name", "test-idp-updated"),
					resource.TestCheckResourceAttr("pangolin_idp_oidc.test", "scopes", "openid profile email"),
				),
			},
		},
	})
}

func testAccIDPOidcConfig() string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_idp_oidc" "test" {
    name        = "test-idp"
    client_id   = "test-client-id"
    client_secret = "test-client-secret"
    auth_url    = "https://example.com/oauth2/auth"
    token_url   = "https://example.com/oauth2/token"
    identifier_path = "$.sub"
    scopes      = "openid"
}
`, testURL, testToken)
}

func testAccIDPOidcConfigUpdated() string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_idp_oidc" "test" {
    name        = "test-idp-updated"
    client_id   = "test-client-id"
    client_secret = "test-client-secret"
    auth_url    = "https://example.com/oauth2/auth"
    token_url   = "https://example.com/oauth2/token"
    identifier_path = "$.sub"
    scopes      = "openid profile email"
}
`, testURL, testToken)
}
