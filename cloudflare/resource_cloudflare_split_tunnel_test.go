package cloudflare

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloudflareSplitTunnel_Include(t *testing.T) {
	// Temporarily unset CLOUDFLARE_API_TOKEN if it is set as the Access
	// service does not yet support the API tokens and it results in
	// misleading state error messages.
	if os.Getenv("CLOUDFLARE_API_TOKEN") != "" {
		defer func(apiToken string) {
			os.Setenv("CLOUDFLARE_API_TOKEN", apiToken)
		}(os.Getenv("CLOUDFLARE_API_TOKEN"))
		os.Setenv("CLOUDFLARE_API_TOKEN", "")
	}

	rnd := generateRandomResourceName()
	name := fmt.Sprintf("cloudflare_split_tunnel.%s", rnd)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccessAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudflareSplitTunnelIncludeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudflareSplitTunnelInclude(rnd, accountID, "example domain", "*.example.com", "include"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "account_id", accountID),
					resource.TestCheckResourceAttr(name, "mode", "include"),
					resource.TestCheckResourceAttr(name, "tunnels.0.description", "example domain"),
					resource.TestCheckResourceAttr(name, "tunnels.0.host", "*.example.com"),
				),
			},
		},
	})
}

func testAccCloudflareSplitTunnelInclude(rnd, accountID string, description string, host string, mode string) string {
	return fmt.Sprintf(`
resource "cloudflare_split_tunnel" "%[1]s" {
	account_id = "%[2]s"
	mode = "%[5]s"
	tunnels {
		description = "%[3]s"
		host = "%[4]s"
	}
}
`, rnd, accountID, description, host, mode)
}

func testAccCheckCloudflareSplitTunnelIncludeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cloudflare.API)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudflare_split_tunnel" {
			continue
		}

		_, err := client.ListSplitTunnels(context.Background(), rs.Primary.Attributes["account_id"], rs.Primary.Attributes["mode"])
		if err == nil {
			return fmt.Errorf("Split Tunnel Include still exists")
		}
	}

	return nil
}

func TestAccCloudflareSplitTunnel_Exclude(t *testing.T) {
	// Temporarily unset CLOUDFLARE_API_TOKEN if it is set as the Access
	// service does not yet support the API tokens and it results in
	// misleading state error messages.
	if os.Getenv("CLOUDFLARE_API_TOKEN") != "" {
		defer func(apiToken string) {
			os.Setenv("CLOUDFLARE_API_TOKEN", apiToken)
		}(os.Getenv("CLOUDFLARE_API_TOKEN"))
		os.Setenv("CLOUDFLARE_API_TOKEN", "")
	}

	rnd := generateRandomResourceName()
	name := fmt.Sprintf("cloudflare_split_tunnel.%s", rnd)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccessAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudflareSplitTunnelExcludeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudflareSplitTunnelExclude(rnd, accountID, "example domain", "*.example.com", "exclude"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "account_id", accountID),
					resource.TestCheckResourceAttr(name, "mode", "exclude"),
					resource.TestCheckResourceAttr(name, "tunnels.0.description", "example domain"),
					resource.TestCheckResourceAttr(name, "tunnels.0.host", "*.example.com"),
				),
			},
		},
	})
}

func testAccCloudflareSplitTunnelExclude(rnd, accountID string, description string, host string, mode string) string {
	return fmt.Sprintf(`
resource "cloudflare_split_tunnel" "%[1]s" {
	account_id = "%[2]s"
	mode = "%[5]s"
	tunnels {
		description= "%[3]s"
		host = "%[4]s"
	}
}
`, rnd, accountID, description, host, mode)
}

func testAccCheckCloudflareSplitTunnelExcludeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*cloudflare.API)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudflare_split_tunnel" {
			continue
		}

		_, err := client.ListSplitTunnels(context.Background(), rs.Primary.Attributes["account_id"], rs.Primary.Attributes["mode"])
		if err == nil {
			return fmt.Errorf("Split Tunnel Exclude still exists")
		}
	}

	return nil
}
