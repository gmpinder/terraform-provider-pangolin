package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTarget_Basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}
	siteID := getTestSiteID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetConfig(siteID, "10.0.0.1", 80),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_target.test", "ip", "10.0.0.1"),
					resource.TestCheckResourceAttr("pangolin_target.test", "port", "80"),
					resource.TestCheckResourceAttrSet("pangolin_target.test", "id"),
				),
			},
			{
				Config: testAccTargetConfig(siteID, "10.0.0.2", 443),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_target.test", "ip", "10.0.0.2"),
					resource.TestCheckResourceAttr("pangolin_target.test", "port", "443"),
				),
			},
		},
	})
}

func testAccTargetConfig(siteID int64, ip string, port int) string {
	return fmt.Sprintf(`
provider "pangolin" {
  base_url = %[1]q
  token    = %[2]q
}

resource "pangolin_resource" "test" {
  org_id    = %[3]q
  name      = "target-test-app"
  protocol  = "tcp"
  http      = true
  subdomain = "target-test"
  domain_id = "local"
}

resource "pangolin_target" "test" {
  resource_id = pangolin_resource.test.id
  site_id     = %[4]d
  ip          = %[5]q
  port        = %[6]d
  enabled     = true
}
`, testURL, testToken, testOrgID, siteID, ip, port)
}
