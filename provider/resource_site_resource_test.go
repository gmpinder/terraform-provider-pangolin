package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSiteResource_Basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}
	siteID := getTestSiteID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig(siteID, "test-app", "host", "app.internal", "app.test-tf.localhost"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "name", "test-app"),
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "mode", "host"),
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "destination", "app.internal"),
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "alias", "app.test-tf.localhost"),
					resource.TestCheckResourceAttrSet("pangolin_site_resource.test", "id"),
				),
			},
			{
				Config: testAccSiteResourceConfig(siteID, "updated-app", "cidr", "10.0.0.0/24", "updated.test-tf.localhost"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "name", "updated-app"),
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "mode", "cidr"),
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "destination", "10.0.0.0/24"),
				),
			},
		},
	})
}

func testAccSiteResourceConfig(siteID int64, name, mode, destination, alias string) string {
	return fmt.Sprintf(`
provider "pangolin" {
  base_url = %[1]q
  token    = %[2]q
}

resource "pangolin_site_resource" "test" {
  org_id      = %[3]q
  site_id     = %[4]d
  name        = %[5]q
  mode        = %[6]q
  destination = %[7]q
  alias       = %[8]q
  enabled     = true
  user_ids    = []
  role_ids    = [1]
  client_ids  = []
  tcp_port_range_string = "*"
  udp_port_range_string = "*"
  disable_icmp          = false
}
`, testURL, testToken, testOrgID, siteID, name, mode, destination, alias)
}
