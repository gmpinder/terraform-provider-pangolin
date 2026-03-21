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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig("test-app", "host", "app.internal", "app.test-tf.localhost"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "name", "test-app"),
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "mode", "host"),
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "destination", "app.internal"),
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "alias", "app.test-tf.localhost"),
					resource.TestCheckResourceAttrSet("pangolin_site_resource.test", "id"),
				),
			},
			{
				Config: testAccSiteResourceConfig("updated-app", "cidr", "10.0.0.0/24", "updated.test-tf.localhost"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "name", "updated-app"),
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "mode", "cidr"),
					resource.TestCheckResourceAttr("pangolin_site_resource.test", "destination", "10.0.0.0/24"),
				),
			},
		},
	})
}

func testAccSiteResourceConfig(name, mode, destination, alias string) string {
	return fmt.Sprintf(`
provider "pangolin" {
  base_url = %[1]q
  token    = %[2]q
}

resource "pangolin_site" "test" {
	name = "test"
	org_id = %[3]q
}

resource "pangolin_site_resource" "test" {
  org_id      = %[3]q
  site_id     = pangolin_site.test.id
  name        = %[4]q
  mode        = %[5]q
  destination = %[6]q
  alias       = %[7]q
  enabled     = true
  user_ids    = []
  role_ids    = [1]
  client_ids  = []
  tcp_port_range_string = "*"
  udp_port_range_string = "*"
  disable_icmp          = false
}

resource "pangolin_site_resource" "test_1" {
  org_id      = %[3]q
  site_id     = pangolin_site.test.id
  name        = %[4]q
  mode        = %[5]q
  destination = %[6]q
  alias       = "other-test.localhost"
}
`, testURL, testToken, testOrgID, name, mode, destination, alias)
}
