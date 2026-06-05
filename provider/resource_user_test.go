package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRole_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRoleConfig("Test Role", "Test Description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_role.test", "name", "Test Role"),
					resource.TestCheckResourceAttr("pangolin_role.test", "description", "Test Description"),
					resource.TestCheckResourceAttrSet("pangolin_role.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccRoleConfig("Updated Role", "Updated Description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pangolin_role.test", "name", "Updated Role"),
					resource.TestCheckResourceAttr("pangolin_role.test", "description", "Updated Description"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRoleConfig(name, description string) string {
	return fmt.Sprintf(`
provider "pangolin" {
  base_url = %[1]q
  token    = %[2]q
}

resource "pangolin_role" "test" {
  org_id      = %[3]q
  name        = %[4]q
  description = %[5]q
}
`, testURL, testToken, testOrgID, name, description)
}
