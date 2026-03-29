package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganization_Basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_organization.test", "name", "test-organization"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "subnet", "100.90.128.0/24"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "utility_subnet", "100.96.128.0/24"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "require_two_factor", "false"),
					resource.TestCheckResourceAttrSet("pangolin_organization.test", "id"),
				),
			},
			{
				Config: testAccOrganizationConfigUpdated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_organization.test", "name", "test-organization-updated"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "require_two_factor", "true"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "max_session_length_hours", "48"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "password_expiry_days", "90"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "settings_log_retention_days_request", "365"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "settings_log_retention_days_access", "365"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "settings_log_retention_days_action", "365"),
				),
			},
		},
	})
}

func TestAccOrganization_WithDefaults(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationConfigMinimal(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_organization.test", "name", "test-org-minimal"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "subnet", "100.90.128.0/24"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "utility_subnet", "100.96.128.0/24"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "require_two_factor", "false"),
				),
			},
			{
				Config: testAccOrganizationConfigWithSpecificDefaults(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_organization.test", "name", "test-org-specific"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "require_two_factor", "true"),
				),
			},
		},
	})
}

func TestAccOrganization_RequiresReplace(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationConfigWithSubnet("100.90.128.0/24"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_organization.test", "name", "test-replace"),
					resource.TestCheckResourceAttr("pangolin_organization.test", "subnet", "100.90.128.0/24"),
				),
			},
			{
				Config:      testAccOrganizationConfigWithSubnet("100.91.128.0/24"),
				ExpectError: regexp.MustCompile(`must be replaced`),
			},
		},
	})
}

func testAccOrganizationConfigBasic() string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_organization" "test" {
    id   = "test-org-id-1"
    name = "test-organization"
}
`, testURL, testToken)
}

func testAccOrganizationConfigUpdated() string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_organization" "test" {
    id                                = "test-org-id-1"
    name                              = "test-organization-updated"
    require_two_factor                = true
    max_session_length_hours          = 48
    password_expiry_days              = 90
    settings_log_retention_days_request = 365
    settings_log_retention_days_access  = 365
    settings_log_retention_days_action  = 365
}
`, testURL, testToken)
}

func testAccOrganizationConfigMinimal() string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_organization" "test" {
    id   = "test-org-id-2"
    name = "test-org-minimal"
}
`, testURL, testToken)
}

func testAccOrganizationConfigWithSpecificDefaults() string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_organization" "test" {
    id   = "test-org-id-2"
    name = "test-org-specific"
    require_two_factor = true
}
`, testURL, testToken)
}

func testAccOrganizationConfigWithSubnet(subnet string) string {
	return fmt.Sprintf(`
provider "pangolin" {
    base_url = %[1]q
    token    = %[2]q
}

resource "pangolin_organization" "test" {
    id   = "test-org-id-3"
    name = "test-replace"
    subnet = %[3]q
}
`, testURL, testToken, subnet)
}
