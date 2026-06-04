package provider

import (
	"testing"

	"github.com/gmpinder/terraform-provider-pangolin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	testOrgID    = "test-tf"
	testToken    = "xt346da9r1xy71z.tg4sgocmc7ecebzaw7acasxct34bbt6duzoaizqs"
	testEmail    = "admin@example.com"
	testPassword = "Passw0rd!"
	testURL      = "http://localhost:3003/v1" // Integration API port and prefix
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"pangolin": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func testAccPreCheck(t *testing.T) {
	// Verify API is reachable
	c := client.NewClient(testURL, testToken)
	_, err := c.ListSites(testOrgID)
	if err != nil {
		t.Fatalf("API unreachable or invalid credentials: %v", err)
	}
}
