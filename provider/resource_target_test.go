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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetConfig("10.0.0.1", 80),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_target.test", "ip", "10.0.0.1"),
					resource.TestCheckResourceAttr("pangolin_target.test", "port", "80"),
					resource.TestCheckResourceAttrSet("pangolin_target.test", "id"),
				),
			},
			{
				Config: testAccTargetConfig("10.0.0.2", 443),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_target.test", "ip", "10.0.0.2"),
					resource.TestCheckResourceAttr("pangolin_target.test", "port", "443"),
				),
			},
		},
	})
}

func testAccTargetConfig(ip string, port int) string {
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

resource "pangolin_site" "test" {
	name = "test"
	org_id = %[3]q
}

resource "pangolin_target" "test" {
	resource_id = pangolin_resource.test.id
	site_id     = pangolin_site.test.id
	ip          = %[4]q
	port        = %[5]d
	enabled     = true
}

resource "pangolin_target" "test_healthcheck" {
	resource_id = pangolin_resource.test.id
	site_id     = pangolin_site.test.id
	ip          = %[4]q
	port        = %[5]d
	enabled     = true

	health_check = {
		enabled  = true
		hostname = %[4]q
		port     = %[5]d
		headers = [
			{
				name  = "Host"
				value = "${pangolin_resource.test.subdomain}"
			}
		]
	}

}
`, testURL, testToken, testOrgID, ip, port)
}
